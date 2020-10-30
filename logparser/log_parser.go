package logparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
)

type LogParser struct {
	// TODO: add mutex for it
	numOfLogLines       int
	numOfParsedLogLines int

	daemonMap    map[string]bool
	daemonLogMap map[string][]*LogEntry // key is daemon name, value is LogEntry array

	serialNumber     string
	firmwareVersions map[string]bool // there can be more than one firmware version in one log file

	parsingDone *atomic.Value
}

func NewLogParser() *LogParser {
	lp := &LogParser{
		numOfLogLines:       0,
		numOfParsedLogLines: 0,
		firmwareVersions:    map[string]bool{},
	}
	lp.parsingDone = &atomic.Value{}
	lp.parsingDone.Store(false)
	return lp
}

func (lp *LogParser) ParsingDone() bool {
	// todo: check ok
	return lp.parsingDone.Load().(bool)
}

func (lp *LogParser) pickCoreContent(logItems []string) (string, bool) {
	res := ""
	isSplunk := false
	for _, logItem := range logItems {
		if len(logItem) == 0 {
			continue
		}
		logItemTmp := strings.Split(logItem, "=")
		if len(logItemTmp) == 2 {
			key := logItemTmp[0]
			value := logItemTmp[1]

			switch key {
			case "time": // we dont record time again
				break
			case "level": // we don't record log level again
				break
			case "#fw":
				if len(value) > 0 && value[0] == '"' {
					value = value[1:]
				}
				if len(value) > 0 && value[len(value)-1] == '"' {
					value = value[:len(value)-1]
				}
				if _, ok := lp.firmwareVersions[value]; !ok {
					lp.firmwareVersions[value] = true
				}
			case "#index": // splunk!
				isSplunk = true
			default: // good stuff, record them
				res += fmt.Sprintf("%s=%s ", key, value)
			}
		} else if logItem[0] != '[' { // we don't record daemon name again
			res += fmt.Sprintf("%s ", logItem)
		}
	}

	// the only returning point
	return res, isSplunk
}

func compressDaemonName(name string) string {
	switch name {
	case "evidencerepaird":
		return "eRepaird"
	case "systembusd":
		return "sysBusd"
	case "qmmf-server":
		return "qmmf"
	default:
		return name
	}
}

func (lp *LogParser) pickTime(logItems []string) string {
	// TODO: regular expression check
	return fmt.Sprintf("%s-%s-%s", logItems[0], logItems[1], logItems[2])
}

func (lp *LogParser) pickLogLevel(logItems []string) LogLevel {
	logLevelItems := strings.Split(logItems[4], ".")
	if len(logLevelItems) == 1 {
		return convertLogLevelString(logLevelItems[0])
	} else if len(logLevelItems) == 2 {
		return convertLogLevelString(logLevelItems[1])
	}
	return LogLevel_Unknown
}
func (lp *LogParser) pickDaemonName(logItems []string) string {
	daemonNameLong := logItems[5]
	daemonNameItems := strings.Split(daemonNameLong, "/")
	daemonName := daemonNameItems[len(daemonNameItems)-1]
	daemonNamePIDIndex := strings.Index(daemonName, "[")
	if daemonNamePIDIndex != -1 {
		daemonName = daemonName[0:daemonNamePIDIndex]
	} else {
		daemonName = daemonName[0 : len(daemonName)-1] // get rid of the last colon
	}
	if len(daemonName) >= 3 && daemonName[0:3] == "ax-" {
		daemonName = daemonName[3:] // get rid of the "ax-" prefix
	}
	return compressDaemonName(daemonName)
}

// Jun 22 01:05:34 X60MMP026 local0.debug ax-gund[3188]: [ax-gund] Status="false" Event="NewEnablingConfigReceived"
func (lp *LogParser) parseOneLogLine(logLine string) (*LogEntry, error) {
	logEntry := &LogEntry{}
	logItems := strings.Fields(logLine)
	if len(logItems) <= 6 {
		return logEntry, fmt.Errorf("not standard log entry: %s", logLine)
	}

	// time
	logEntry.TimeString = lp.pickTime(logItems)

	// serial number
	lp.serialNumber = logItems[3]

	// log level
	logEntry.LogLevel = lp.pickLogLevel(logItems)

	// daemon name
	logEntry.DaemonName = lp.pickDaemonName(logItems)

	// core content
	isSplunk := false
	logEntry.CoreContent, isSplunk = lp.pickCoreContent(logItems[6:])
	// splunk?
	logEntry.IsSplunk = isSplunk

	return logEntry, nil
}

func (lp *LogParser) reset() {
	lp.parsingDone.Store(false)
	lp.numOfParsedLogLines = 0
	lp.numOfLogLines = 0
	lp.daemonMap = map[string]bool{}
	lp.daemonLogMap = map[string][]*LogEntry{}
	lp.firmwareVersions = map[string]bool{}
}

// ParseLogString take a string value (the log) as input then return an array of log entry,
// each log entry is parsed from each line in input
func (lp *LogParser) ParseLogString(reader io.Reader) ([]*LogEntry, error) {
	lp.reset()
	// split log by lines
	logEntries := make([]*LogEntry, 0)
	lp.numOfLogLines = int(0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lp.numOfLogLines++
		logEntry, err := lp.parseOneLogLine(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			logEntry.AbsoluteIndex = lp.numOfLogLines - 1
			logEntries = append(logEntries, logEntry)
			// record daemon dict and daemon to logentries dict
			if _, ok := lp.daemonMap[logEntry.DaemonName]; !ok {
				lp.daemonMap[logEntry.DaemonName] = true
				lp.daemonLogMap[logEntry.DaemonName] = []*LogEntry{logEntry}
			} else {
				lp.daemonLogMap[logEntry.DaemonName] = append(lp.daemonLogMap[logEntry.DaemonName], logEntry)
			}
		}
		lp.numOfParsedLogLines++
	}

	lp.parsingDone.Store(true)
	fmt.Printf("Parsing finished, in total %d (vs %d) lines of valid log\n", len(logEntries), lp.numOfLogLines)

	return logEntries, nil
}

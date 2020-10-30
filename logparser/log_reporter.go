package logparser

import (
	"fmt"
	"sort"
)

// ReportParsingProgress returns progress of parsing. If no parsing is on going it returns
// previous parsing task percentage. If no previous parsing exists, it will return 0
func (lp *LogParser) ReportParsingProgress() int {
	if lp.numOfLogLines == 0 {
		return 0
	}

	if lp.ParsingDone() {
		return 100
	}

	return 0
}

// GetDaemonList returns an array of string and each entry is the name of daemon, without emoji
func (lp *LogParser) GetDaemonList() []string {
	daemonArray := make([]string, len(lp.daemonMap))
	index := 0
	for key := range lp.daemonMap {
		daemonArray[index] = key
		index++
	}
	sort.Slice(daemonArray, func(i, j int) bool {
		return daemonArray[i] < daemonArray[j]
	})

	fmt.Println("daemon list fetched")

	return daemonArray
}

// GetDaemonListWithEmoji returns an array of string and each entry is the name of daemon, with emoji
func (lp *LogParser) GetDaemonListWithEmoji() []string {
	return []string{}
}

// GetLogByDaemonName returns a string containing log about a daemon, separated by newline
func (lp *LogParser) GetLogByDaemonName(daemonName string) []*LogEntry {
	logEntries, ok := lp.daemonLogMap[daemonName]
	if !ok {
		return []*LogEntry{}
	}

	return logEntries
}

// LogMetadata represents the metadata of this log file, like serial number
type LogMetadata struct {
	FirmwareVersions     []string
	SerialNumber         string
	NumberOfDaemonLogged int
	NumberOfLogEntry     int64
}

// GetMetaData returns metadata about this log file, defined in LogMetadata
func (lp *LogParser) GetMetaData() *LogMetadata {
	FirmwareVersions := []string{}
	for key := range lp.firmwareVersions {
		FirmwareVersions = append(FirmwareVersions, key)
	}
	return &LogMetadata{
		FirmwareVersions:     FirmwareVersions,
		SerialNumber:         lp.serialNumber,
		NumberOfDaemonLogged: len(lp.daemonMap),
		NumberOfLogEntry:     int64(lp.numOfParsedLogLines),
	}
}

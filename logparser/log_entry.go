package logparser

import "fmt"

type LogLevel int

const (
	LogLevel_Debug   LogLevel = 0
	LogLevel_Info    LogLevel = 1
	LogLevel_Warn    LogLevel = 2
	LogLevel_Error   LogLevel = 3
	LogLevel_Crit    LogLevel = 4
	LogLevel_Unknown LogLevel = 5
)

func convertLogLevelString(logLevelString string) LogLevel {
	switch logLevelString {
	case "debug":
		return LogLevel_Debug
	case "warn":
		return LogLevel_Warn
	case "info":
		return LogLevel_Info
	case "err":
		return LogLevel_Error
	case "crit":
		return LogLevel_Crit
	default:
		return LogLevel_Unknown
	}
}

type LogEntry struct {
	TimeString    string
	LogLevel      LogLevel
	DaemonName    string
	CoreContent   string
	IsSplunk      bool
	AbsoluteIndex int // the absolute index in engineering log with all other daemons considered
}

func printLogEntry(entry *LogEntry) {
	fmt.Println(entry)
}

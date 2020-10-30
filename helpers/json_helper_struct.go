// This file contains helper struct for sending struct from golang to javascipt by JSON.marshal
package helpers

import (
	"git.taservs.net/Metropolis/SupermanLogViewer/SupermanLogViewer-backend/logparser"
)

// DaemonList is used for Marshaling to json then sent to frontend
type DaemonList struct {
	List []string
}

// LogEntryList list of *LogEntry
type LogEntryList struct {
	List []*logparser.LogEntry
}

type SplitViewEntry struct {
	LogEntry    *logparser.LogEntry
	ColumnIndex int
}

type SplitViewLogEntryList struct {
	List []*SplitViewEntry
}

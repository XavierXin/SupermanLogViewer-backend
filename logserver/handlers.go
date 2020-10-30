package logserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"git.taservs.net/Metropolis/SupermanLogViewer/SupermanLogViewer-backend/helpers"
)

// CORS enables sending/getting request from different origin.
// It was disabled due to security concerns
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (ls *LogServer) handleLogParsing(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	_, err := ls.logParser.ParseLogString(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ls *LogServer) handleParsingProgress(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Fprintf(w, strconv.Itoa(ls.logParser.ReportParsingProgress()))
}

func (ls *LogServer) handleDaemonList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	daemonList := ls.logParser.GetDaemonList()
	jsonMarshaledDaemonList, err := json.Marshal(helpers.DaemonList{
		List: daemonList,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, string(jsonMarshaledDaemonList))
}

func (ls *LogServer) handleMetadata(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	metadata := ls.logParser.GetMetaData()
	jsonMarshaledMetadata, err := json.Marshal(metadata)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, string(jsonMarshaledMetadata))
}

func (ls *LogServer) handleDaemonLog(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	daemonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	logEntries := ls.logParser.GetLogByDaemonName(string(daemonBytes))
	jsonMarshaledLogEntries, err := json.Marshal(&helpers.LogEntryList{
		List: logEntries,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprintf(w, string(jsonMarshaledLogEntries))
}

func (ls *LogServer) handleMultiDaemonView(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jsonDaemonList, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var daemonList = struct {
		List []string
	}{
		List: []string{},
	}

	if err = json.Unmarshal(jsonDaemonList, &daemonList); err != nil {
		fmt.Println(err.Error())
	}

	// daemonLogMap[i][j] is the jth log entry of ith daemon in daemonList
	daemonLogList := make([]*helpers.LogEntryList, len(daemonList.List))

	for i, daemonName := range daemonList.List {
		logEntries := ls.logParser.GetLogByDaemonName(daemonName)
		daemonLogList[i] = &helpers.LogEntryList{
			List: logEntries,
		}
	}

	splitDaemonView := helpers.GenerateSplitViewList(daemonLogList)
	jsonMarshaledSplitViewList, err := json.Marshal(&helpers.SplitViewLogEntryList{
		List: splitDaemonView,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprintf(w, string(jsonMarshaledSplitViewList))

	fmt.Println("split view generated and sent")
}

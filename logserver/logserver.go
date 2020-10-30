package logserver

import (
	"context"
	"log"
	"net/http"

	"git.taservs.net/Metropolis/SupermanLogViewer/SupermanLogViewer-backend/logparser"
)

type LogServer struct {
	logParser  *logparser.LogParser
	httpServer *http.Server
}

func NewLogServer() *LogServer {
	return &LogServer{
		logParser: logparser.NewLogParser(),
	}
}

// StartBlock starts the Golang log parsing server. This function will block.
func (ls *LogServer) StartBlock() {
	ls.httpServer = &http.Server{Addr: "localhost:7777"}

	// Serve /logparser with a text response.
	http.HandleFunc("/logparser", ls.handleLogParsing)

	// return parsing progress, returned value looks like "99", "100" means finished
	http.HandleFunc("/parsingProgress", ls.handleParsingProgress)

	http.HandleFunc("/askForDaemonList", ls.handleDaemonList)

	http.HandleFunc("/askForMetadata", ls.handleMetadata)

	http.HandleFunc("/askForDaemonLog", ls.handleDaemonLog)

	http.HandleFunc("/askForMultiDaemonView", ls.handleMultiDaemonView)

	// Start the server at http://localhost:port
	if err := ls.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (ls *LogServer) Stop() {
	if err := ls.httpServer.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

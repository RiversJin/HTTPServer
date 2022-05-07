package main

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

type fileHandlerWithLog struct {
	handler http.Handler
	logger  func(*http.Request)
}

func FileHandlerWithLog(root http.FileSystem, logger func(*http.Request)) http.Handler {
	return &fileHandlerWithLog{gziphandler.GzipHandler(http.FileServer(root)), logger}
}

func (fhw *fileHandlerWithLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fhw.logger(r)

	fhw.handler.ServeHTTP(w, r)
}

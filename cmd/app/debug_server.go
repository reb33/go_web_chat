package main

import (
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
)

const (
	debugServerHost = "127.0.0.1"
	debugServerReadHeaderTimeout = 5 * time.Second
)

func newDebugMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}

func newDebugServer(port int) *http.Server {
	return &http.Server{
		Addr:              net.JoinHostPort(debugServerHost, strconv.Itoa(port)),
		Handler:           newDebugMux(),
		ReadHeaderTimeout: debugServerReadHeaderTimeout,
	}
}

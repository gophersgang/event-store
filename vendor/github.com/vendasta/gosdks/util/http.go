package util

import (
	"fmt"
	"net/http"
	"net"
	"net/http/pprof"

	"golang.org/x/net/context"

    "github.com/vendasta/gosdks/logging"
)


// StartServer bootstraps the http server for this instance (blocking)
func StartHTTPServer(healthz http.HandlerFunc, port int, mux *http.ServeMux) (error) {
	httpSrv := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler: mux,
	}

    var lis net.Listener
    var err error
	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
        logging.Errorf(context.Background(), "Error creating HTTP listening socket: %s", err.Error())
		return err
	}
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
    mux.Handle("/healthz", healthz)

	return httpSrv.Serve(lis)
}

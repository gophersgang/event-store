package main

import (
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/vendasta/gosdks/config"
	"github.com/vendasta/gosdks/config/elastic"
	"github.com/vendasta/gosdks/logging"
	"github.com/vendasta/gosdks/statsd"
	"github.com/vendasta/gosdks/util"
	"github.com/vendasta/event-store/pb"
	"github.com/vendasta/event-store/pkg/api"
)

const (
	APP_NAME  = "local"
	GRPC_PORT = 11000
	HTTP_PORT = 11001
)

func main() {
	var err error
	ctx := context.Background()

	//Setup Application logging and switch the logger
	if !config.IsLocal() {
		namespace := config.GetGkeNamespace()
		podName := config.GetGkePodName()
		if err = logging.Initialize(namespace, podName, APP_NAME); err != nil {
			logging.Errorf(ctx, "Error initializing logger: %s", err.Error())
			os.Exit(-1)
		}
	}

	//Setup ElasticSearch Client
	if err = elasticclient.Initialize(); err != nil {
		logging.Errorf(ctx, "Error initilizing Elastic Client: %s", err.Error())
		os.Exit(-1)
	}

	//Setup StatsD Client
	if err = statsd.Initialize(APP_NAME); err != nil {
		logging.Errorf(ctx, "Error initilizing statsd client: %s", err.Error())
		os.Exit(-1)
	}

	//TODO: (optional) INSERT YOUR CUSTOM AUTHORIZATION INTERCEPTOR HERE
	var identityInterceptor = util.NoAuthInterceptor

	//Create Logging Interceptor
	var loggingInterceptor = logging.Interceptor()

	//Create a GRPC Server
	logging.Infof(ctx, "Creating GRPC server...")
	grpcServer := util.CreateGrpcServer(loggingInterceptor, identityInterceptor)

	//--------- INSERT YOUR CODE HERE ------------
	pb.RegisterEventStoreServer(grpcServer, &api.EventStoreServer{})
	//REGISTER_GRPC_SERVERS_HERE

	//Start GRPC API Server
	go func() {
		logging.Infof(ctx, "Running GRPC server...")
		if err = util.StartGrpcServer(grpcServer, GRPC_PORT); err != nil {
			logging.Errorf(ctx, "Error starting GRPC Server: %s", err.Error())
			os.Exit(-1)
		}
	}()

	//Start Healthz and Debug HTTP API Server
	healthz := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//TODO: (optional) INSERT YOUR CODE HERE
		return
	}

	logging.Infof(ctx, "Running HTTP server...")
	mux := http.NewServeMux()
	if err = util.StartHTTPServer(healthz, HTTP_PORT, mux); err != nil {
		logging.Errorf(ctx, "Error starting Healthz & Debug server: %s", err.Error())
		os.Exit(-1)
	}
}

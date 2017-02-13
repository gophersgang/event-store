# MSCLI

This is a command line interface intended to be used for Vendasta Microservices.

All the instructions below __assume__ you are following the "Single Go Path" pattern.

# Projects in the Single GOPATH
This recommends we use a single go path approach to our project structure. This 
means all our projects would fall under `$HOME/go`. Your folders would look like:

```
/User/<USERNAME>/go/src/github.com/vendasta/AA
/User/<USERNAME>/go/src/github.com/vendasta/ARM
/User/<USERNAME>/go/src/github.com/vendasta/BS
...
```

To set this up use the following steps:

1. Install golang from `https://golang.org/dl/`
1. Make the project folder `mkdir -p ~/go/src/github.com/vendasta/`
1. Add the following to your `~/.bash_profile`
  1. `export GOPATH=$HOME/go`
  1. `export GOBIN=$GOPATH/bin`
  1. `export PATH=$GOBIN:$PATH`
1. Source the `~/.bash_profile` (ie restart terminal or whatever)

# Installation Instructions

## Get the mscli
Later the binary might be distributed, for now you have to build it yourself. 
You can either use `go get` or a `git clone`. If you want to contribute back to a repo you should clone it (this will include the .git files needed). If you just want to use the repo you can run `go get` (this will not include the .git files).

### Git clone
1. `cd ~/go/src/github.com/vendasta/`
1. `git clone https://github.com/vendasta/gosdks.git`

### Go Get
`go get github.com/vendasta/gosdks/tools/mscli`

### Installation
`go install github.com/vendasta/gosdks/tools/mscli`

### AND THEN
1. Install [glide](http://glide.sh/).
1. Install [minikube](https://github.com/kubernetes/minikube/releases).
1. Install goreturns `go get github.com/sqs/goreturns`

# Getting a Microservice Started
Run `mscli` to see a list of commands. The one you are going to use to setup your microservice is setupAll

## Create a new Microservice in an old project
Because this will probably be a common path for the creation of microservices we thought we'd expand some on some specifics.
Currently this means you will probably have duplicate copy of any 'legacy' projects if you are adding microservices to them. This will be changed in the future.

1. `cd ~/go/src/github.com/vendasta/`
1. If you haven't added the project to the 1 go-path directory `git clone <clone link for product>`
1. `cd <product name>`
1. Branch. (`git checkout -b <your branch name>`)
1. `mkdir <your new microservice name>`
1. `cd <microservice name>`

## Define your proto file
Here's an example
```proto
syntax = "proto3";

package pb;

import "google/protobuf/empty.proto";

message PubsubMetric{
    int64 topics = 1;
    int64 subscriptions = 2;
}

message GetPubsubMetricsResponse {
    PubsubMetric pubsubMetric = 1;
}

service Metrics {
    rpc Get(google.protobuf.Empty) returns (GetPubsubMetricsResponse);
}
```
This will end up with 1 endpoint that returns the count of topics and subscriptions

There's a __public__ repo for storing protos that are intended to be exposed directly to customers at [here](https://github.com/vendasta/vendastaapis). 
If you are writing a proto for internal use, you might want to define it somewhere else.

Once your protos are defined:

1. `cd ~/go/src/github.com/vendasta/<project name>/<microservice name>`
1. `mscli setupAll --path=<path to proto directory containing proto files>`

You will now have some files generated
* microservice.yaml
  * This has the configuration for deploying your microservice. It should be generated with the environments local and test. 
You may want to fill out extra values in this before building/deploying
* pb/
  * This directory contains the generated protobuf files
* server/
  * main.go
    * The entry point for your application. If you've generated your protos, those services should already be registered in here
    * You probably don't need to touch this for a basic setup
  * [proto file name].go
    * If you've generated your protos you will have one generated file for each proto file. This is where you should add your implementation

## Building a stub GRPC Service
Let's say you named your proto file metrics.proto. After running setupAll, you will now have a `server` directory with the files `main.go` and `metrics.go`.
Open `metrics.go` and add your implementation.

## Building and Deploying your container (locally)

1. `minikube start`
1. Make sure you're building containers that work for minikube instead of Docker for Mac: `eval $(minikube docker-env)`
1. `cd <path to where microservice.yaml is located>`
1. `mscli build <version>`
1. `mscli deploy <version>`
1. Edit your /etc/hosts file to include an entry "192.168.99.100 minikube.xxx".  Now you can get to your minikube installation using the hostname minikube.xxx (google requires a TLD like xxx to support oauth2, that's why we use .xxx)
1. View the health of your cluster [http://minikube.xxx:30000](http://minikube.xxx:30000)
1. View the debug server [https://minikube.xxx:31956/debug/pprof](https://minikube.xxx:31956/debug/pprof)
1. Your GRPC service should now be serving at [minikube.xxx:31957]



## Testing your GRPC Service
Here is a sample client to see if your grpc service is running. Drop it in to `/path to where microservice.yaml is/client/main.go`
You will also need to change this import `github.com/vendasta/rd-metrics/pubsub/pb` to be `github.com/vendasta/where your microservice.yaml is/pb`
```golang
package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/vendasta/rd-metrics/pubsub/pb"
)

func main() {
	conn, err := grpc.Dial("minikube.xxx:31957", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMetricsClient(conn)

	r, err := c.Get(context.Background(), &google_protobuf.Empty{})
	if err != nil {
		log.Fatalf("could not get stats: %v", err)
	}

	log.Printf("Topics %d", r.PubsubMetric.Topics)
	log.Printf("Subscriptions %d", r.PubsubMetric.Subscriptions)
}

```

## WARNING
`mscli protoc` and `mscli setupAll` are both destructive operations for the directories `/path to microservice.yaml/pb` and `/path to microservice.yaml/server/`




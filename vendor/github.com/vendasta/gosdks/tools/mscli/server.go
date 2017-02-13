package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const (
	gRPCServerRegisterTag = "//REGISTER_GRPC_SERVERS_HERE"
	SERVER_TEMPLATE       = `package main

import (
    "net/http"
    "os"

    "golang.org/x/net/context"
    
	"github.com/vendasta/gosdks/logging"
    "github.com/vendasta/gosdks/config"
    "github.com/vendasta/gosdks/config/elastic"
    "github.com/vendasta/gosdks/statsd"
    "github.com/vendasta/gosdks/util"
)

const (
    APP_NAME = "{{.Name}}" 
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
            logging.Errorf(ctx, "Error initializing logger: %s", err.Error());
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

    //Create a vStore client
    logging.Infof(ctx, "Creating vStore Client...")
	vstoreClient, err := vstore.New()
	if err != nil {
		logging.Errorf(ctx, "Error initializing vstore client %s", err.Error())
		os.Exit(-1)
	}
	logging.Infof(ctx, "Using vStore Client: %#v", vstoreClient)

    //--------- INSERT YOUR CODE HERE ------------
	` + gRPCServerRegisterTag + `

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
`

	GLIDE_TEMPLATE = `package: {{.GoPackageName}}/server
import:
- package: github.com/golang/protobuf
  subpackages:
  - proto
  - ptypes/empty
  - ptypes/timestamp

- package: github.com/vendasta/gosdks
  subpackages:
  - config
  - config/elastic
  - logging
  - util

- package: golang.org/x/net
  subpackages:
  - context

- package: google.golang.org/grpc
  version: ^1.0.4

- package: google.golang.org/genproto
  version: master

- package: cloud.google.com/go

- package: google.golang.org/api
`
)

func CreateGrpcBoilerplate(config MicroserviceConfig) {
	fmt.Printf("Creating GRPC boilerplate\n")
	var tmpl *template.Template
	var err error

	if tmpl, err = template.New("microservice_golang_boiler").Parse(SERVER_TEMPLATE); err != nil {
		log.Fatalf("Error creating golang temnplate: %s", err.Error())
	}

	buf := bytes.NewBufferString("")
	envConfig := config.GetEnvironment()
	if err = tmpl.Execute(buf, envConfig); err != nil {
		log.Fatalf("Error executing golang template: %s", err.Error())
	}

	if config.Debug {
		fmt.Printf("------- Golang server: main.go --------\n")
		fmt.Printf("%s", buf.String())
		fmt.Printf("---------------------------------------\n")
	}

	var f *os.File
	if f, err = os.Create("./server/main.go"); err != nil {
		log.Fatalf("Error creating golang server file main.go: %s", err.Error())
	}
	defer f.Close()
	if _, err = f.WriteString(buf.String()); err != nil {
		log.Fatalf("Error writing golang server file main.go: %s", err.Error())
	}
}

func CreateGlideYaml(config MicroserviceConfig) {
	fmt.Printf("Creating glide.yaml\n")
	var tmpl *template.Template
	var err error

	if tmpl, err = template.New("microservice_glide_yaml_boiler").Parse(GLIDE_TEMPLATE); err != nil {
		log.Fatalf("Error creating glide.yaml temnplate: %s", err.Error())
	}

	buf := bytes.NewBufferString("")
	//envConfig := config.GetEnvironment()
	if err = tmpl.Execute(buf, config); err != nil {
		log.Fatalf("Error executing glide.yaml template: %s", err.Error())
	}

	if config.Debug {
		fmt.Printf("------- Golang server: glide.yaml --------\n")
		fmt.Printf("%s", buf.String())
		fmt.Printf("---------------------------------------\n")
	}

	var f *os.File
	if f, err = os.Create("./glide.yaml"); err != nil {
		log.Fatalf("Error creating glide.yaml: %s", err.Error())
	}
	defer f.Close()
	if _, err = f.WriteString(buf.String()); err != nil {
		log.Fatalf("Error writing golang server file glide.yaml: %s", err.Error())
	}
}

func RunGlide() {
	fmt.Println("Running glide commands `glide cc` and `glide up -v` this will take a while")
	cmd := exec.Command("glide", "cc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\n Error running command: %v, %v\n", cmd, err)
	}
	cmd = exec.Command("glide", "up", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\n Error running command: %v, %v\n", cmd, err)
	}
}

func BuildSource(config MicroserviceConfig) {
	var err error
	cmd := exec.Command("go", "build", "-o", "./server/server", config.GoPackageName)
	env := os.Environ()
	env = append(env, "GOOS=linux")
	env = append(env, "GOARCH=amd64")
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("Error starting build: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors building: %s", err.Error())
	}
}

func RunTests(config MicroserviceConfig, dockerTag string) {
	fmt.Printf("Running Tests\n")
	var err error
	cmd := exec.Command("docker", "run", "--env", "ENVIRONMENT=test", dockerTag, "go", "test", fmt.Sprintf("../%s/pkg/...", config.GoPackageName))
	env := os.Environ()
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("Error running tests: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors testing: %s", err.Error())
	}
}

func RunLint(config MicroserviceConfig, dockerTag string) {
	fmt.Printf("Linting\n")
	var err error
	cmd := exec.Command("docker", "run", dockerTag, "golint", "-set_exit_status", fmt.Sprintf("/go/src/%s/pkg/...", config.GoPackageName))
	fmt.Printf("Running %s %s\n", cmd.Path, cmd.Args)
	env := os.Environ()
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("Error running linting: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors linting: %s", err.Error())
	}
}

func RunVet(config MicroserviceConfig, dockerTag string) {
	fmt.Printf("Vetting\n")
	var err error
	cmd := exec.Command("docker", "run", dockerTag, "go", "tool", "vet", fmt.Sprintf("/go/src/%s/pkg", config.GoPackageName))
	fmt.Printf("Running %s %s\n", cmd.Path, cmd.Args)
	env := os.Environ()
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("Error running linting: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors linting: %s", err.Error())
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"os/exec"

	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/rest"
	"time"
	"strconv"
)

const (
	pathToServerDir = "./server"
	pathToPkgDir = "./pkg"
	pathToAPIDir = pathToPkgDir + "/api"
)

var (
	app = kingpin.New("mscli", "A Vendasta microservice cli tool")
	debug = app.Flag("debug", "Show verbose information").Bool()
	settings = app.Flag("settings", "YAML file with microservice settings").Default("microservice.yaml").String()
	environment = app.Flag("env", "Environment to target").Default("local").String()

	buildConfig = app.Command("gen-config", "Creates a basic configuration")

	bootstrap = app.Command("bootstrap", "Creates a new microservice")

	build = app.Command("build", "Builds a new Docker Image (locally)")
	buildVersion = build.Arg("version", "The version to attach to the docker image").Required().String()
	buildLocation = build.Flag("location", "Where to build the docker image (local, cloud)").Default("local").String()

	deploy = app.Command("deploy", "Deploys a service")
	deployVersion = deploy.Arg("version", "The version of the docker image to run").String()

	setupAll = app.Command("setupAll", "Runs all the steps required for generating a microservice from proto file(s)")
	setupAllPath = setupAll.Flag("path", "Path to the the directory that contains 1 or more proto files").Required().String()

	test = app.Command("test", "Runs tests on your microservice.")
	testVersion = test.Arg("version", "The version of the docker image to run").Required().String()

	lint = app.Command("lint", "Runs lint on your microservice.")
	lintVersion = lint.Arg("version", "The version of the docker image to run").Required().String()

	vet = app.Command("vet", "Runs go vet on your microservice.")
	vetVersion = vet.Arg("version", "The version of the docker image to run").Required().String()

	protoc = app.Command("protoc", "Compile and build protos")
	protocPath = protoc.Flag("path", "Path to the the directory that contains 1 or more proto files").Required().String()

	generateJWT = app.Command("jwt", "Generates a JWT with vendasta local creds to call a local microservice.")
	scope = generateJWT.Arg("scope", "Scope is specified in the auth section for the endpoints audience.").Required().String()
	version = app.Command("version", "Returns the version of the cli")
)

func main() {
	//if !inCluster() {
	//	checkDependancies()
	//}
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if command == setupAll.FullCommand() {
		fmt.Printf("Sit back, relax, and watch this\n")
		runCommand(buildConfig.FullCommand())
		runCommand(bootstrap.FullCommand())
		protocPath = setupAllPath
		runCommand(protoc.FullCommand())
	} else {
		runCommand(command)
	}

}

func runCommand(command string) {
	switch command {
	case buildConfig.FullCommand():
		fmt.Printf("Generating a shell microservice.yaml file\n")

		cf := NewConfigFile()
		FailIfFileExists("./microservice.yaml")
		WriteConfigFile(cf, "./microservice.yaml")

	case bootstrap.FullCommand():
		fmt.Printf("Bootstrapping a new microservice\n")
		EnsureDirExists(pathToServerDir)

		//Create a configuration
		config := parseConfig(*settings)
		config.Debug = *debug
		config.Environment = *environment

		CreateGrpcBoilerplate(config)
		CreateContinuousIntegrationBoilerplate(config)
		CreateGlideYaml(config)
		RunGlide()
		fmt.Printf("Completed Bootstrapping.\n")

	case build.FullCommand():
		fmt.Printf("Building a new Docker Image\n")

		//Create a configuration
		config := parseConfig(*settings)
		config.Debug = *debug

		//BuildSource(config)
		tag := BuildDockerImage(config, *buildVersion)
		fmt.Printf("Docker image was built with tag %s\n", tag)

	//BuildDockerImage(config, *buildVersion, *buildLocation)

	case deploy.FullCommand():
		// if a version has not been specified, generate one based on a timestamp
		v := *deployVersion
		if v == "" {
			v = strconv.FormatInt(time.Now().UTC().Unix(), 10)
		}
		fmt.Printf("Deploying to Kubernetes. Version: %s\n", v)

		//Create a configuration
		config := parseConfig(*settings)
		config.Debug = *debug
		config.Environment = *environment
		config.Version = v

		//Get a Kubernetes ClientSet
		clientSet := GetK8sClientSet(config)

		//BuildSource(config)
		dockerImage := BuildDockerImage(config, v)

		//Copy the image to minikube if necessary
		if config.Environment == "local" {
			//CopyToMinikube(config)
		} else {
			PushToContainerRegistry(dockerImage)
		}

		//Create the configMaps
		CreateConfigMap(config, clientSet)

		//Create the services
		CreateServices(config, clientSet)

		//Create/update the deployment
		CreateDeployment(config, clientSet)

		//Create side apps if needed
		CreateAppsDeployment(config, clientSet)

		//Create/update the secret
		CreateSecret(config, clientSet)

		// Create/Update the horizontal pod autoscaler.
		CreateHorizontalPodAutoscaler(config, clientSet)

		//Create/update the local proxy deployment
		if config.GetEnvironment().Name == "local" {
			CreateLocalProxyDeployment(clientSet)
			CreateLocalProxyService(clientSet)
		}

	case protoc.FullCommand():
		fmt.Printf("%s\n", "Compiling protos")
		EnsureDirExists(pathToServerDir)
		EnsureDirExists(pathToPkgDir)
		EnsureDirExists(pathToAPIDir)

		config := parseConfig(*settings)
		config.Debug = *debug
		BuildFromProtos(config, *protocPath)

	case test.FullCommand():
		config := parseConfig(*settings)
		tag := BuildDockerImage(config, *testVersion)
		RunTests(config, tag)

	case lint.FullCommand():
		config := parseConfig(*settings)
		tag := BuildDockerImage(config, *lintVersion)
		RunLint(config, tag)

	case vet.FullCommand():
		config := parseConfig(*settings)
		tag := BuildDockerImage(config, *vetVersion)
		RunVet(config, tag)

	case generateJWT.FullCommand():
		GenerateJWT(*scope)

	case version.FullCommand():
		fmt.Printf("%s\n", CLI_VERSION)

	default:
		fmt.Printf("Unknown command\n")
	}
}

type dependency struct {
	program    string
	installURL string
}

func checkDependancies() {

	dependencies := []dependency{
		dependency{program: "go", installURL: "https://golang.org/doc/install"},
		dependency{program: "minikube", installURL: "https://github.com/kubernetes/minikube/releases"},
		dependency{program: "glide", installURL: "http://glide.sh/"},
		dependency{program: "goreturns", installURL: "https://github.com/sqs/goreturns"},
	}

	success := true
	for _, dep := range dependencies {
		checkInstall := exec.Command("which", fmt.Sprintf("%s", dep.program))
		if err := checkInstall.Run(); err != nil {
			success = false
			fmt.Printf("\033[0;31m %s not installed, go install it now: %s \033[0m \n", dep.program, dep.installURL)
		}
	}

	if !success {
		log.Fatal("Some dependancies are not installed. Please install them before continuing.")
	}
}
func inCluster() bool {
	_, err := rest.InClusterConfig(); if err == nil {
		return true
	}
	return false
}

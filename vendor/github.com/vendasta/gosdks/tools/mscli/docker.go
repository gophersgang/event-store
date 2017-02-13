package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
	"strings"
)

const DOCKERFILE_TEMPLATE = `FROM golang:1.7
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN go get -u github.com/golang/lint/golint

COPY . /go/src/{{.GoPackageName }}

RUN go build -o /bin/{{.Name}} {{.GoPackageName }}/server/
CMD ["/bin/{{.Name}}"]
`

func CreateDockerfileBoilerplate(config MicroserviceConfig) {
	var tmpl *template.Template
	var err error

	if tmpl, err = template.New("microservice_golang_boiler").Parse(DOCKERFILE_TEMPLATE); err != nil {
		log.Fatalf("Error creating golang temnplate: %s", err.Error())
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, config); err != nil {
		log.Fatalf("Error executing golang template: %s", err.Error())
	}

	if config.Debug {
		fmt.Printf("------- Golang server: Dockerfile --------\n")
		fmt.Printf("%s", buf.String())
		fmt.Printf("---------------------------------------\n")
	}

	var f *os.File
	if f, err = os.Create("./Dockerfile"); err != nil {
		log.Fatalf("Error creating Dockerfile: %s", err.Error())
	}
	defer f.Close()
	if _, err = f.WriteString(buf.String()); err != nil {
		log.Fatalf("Error writing Dockerfile: %s", err.Error())
	}
}

func BuildDockerImage(config MicroserviceConfig, version string) string {
	fmt.Printf("Building docker image\n")
	//Create the Dockerfile
	CreateDockerfileBoilerplate(config)

	//Build the Docker image
	var err error
	tag := fmt.Sprintf("gcr.io/repcore-prod/%s:%s", config.Name, version)
	cmd := exec.Command("docker", "build", "-t", tag, ".")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if !IsOnJenkins() {
		env := os.Environ()
		env = append(env, GetMinikubeEnv()...)
		cmd.Env = env
		log.Printf("Using env: %s to build image", cmd.Env)
	}

	if err = cmd.Start(); err != nil {
		log.Fatalf("Error starting docker build: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors building docker image: %s", err.Error())
	}

	//Delete the Dockerfile
	if err = os.Remove("./Dockerfile"); err != nil {
		log.Fatalf("Error removing Dockerfile: %s", err.Error())
	}
	return tag
}

func CopyToMinikube(config MicroserviceConfig, version string) {
	//Save the image to a .tar File
	var err error
	tag := fmt.Sprintf("gcr.io/repcore-prod/%s:%s", config.Name, version)
	cmd := exec.Command("docker", "save", "-o", "temp_image", tag)
	cmd.Dir = "./server"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("There were errors starting saving the local docker image: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors saving the local docker image: %s", err.Error())
	}

	//Get the minikube environment
	minikubeEnv := GetMinikubeEnv()

	//Import the docker image into minikube
	cmd = exec.Command("docker", "load", "-i", "temp_image")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env := os.Environ()
	env = append(env, minikubeEnv...)
	cmd.Env = env

	if err = cmd.Start(); err != nil {
		log.Fatalf("There were errors starting loading the local docker image: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors loading the local docker image: %s", err.Error())
	}

	//Delete the image
	if err = os.Remove("./Dockerfile"); err != nil {
		log.Fatalf("Error removing image file: %s", err.Error())
	}
}

func PushToContainerRegistry(dockerImage string) {
	dockerLogin()

	cmd := exec.Command("docker", "push", dockerImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env := os.Environ()
	cmd.Env = env
	err := cmd.Start()
	if err != nil {
		log.Fatalf("There were errors pushing the image to gcr. %s", err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("There were errors pushing the image to gcr. %s", err.Error())
	}
}

func GetMinikubeEnv() []string {
	var stdout io.ReadCloser
	cmd := exec.Command("minikube", "docker-env")
	stdout, err := cmd.StdoutPipe(); if err != nil {
		log.Fatalf("Error collecting minikube env %s", err.Error())
	}
	cmd.Start()

	r := bufio.NewReader(stdout)
	o, err := ioutil.ReadAll(r); if err != nil {
		log.Fatalf("Error collecting minikube env %s", err.Error())
	}
	toReturn := []string{}
	for _, line := range strings.Split(string(o[:]), "\n") {
		if strings.HasPrefix(line, "export") {
			toReturn = append(toReturn, strings.Replace(strings.TrimSpace(strings.Split(line, "export")[1]), "\"", "", -1))
		}
	}
	return toReturn
}

func IsOnJenkins() bool {
	return os.Getenv("JENKINS_HOME") != ""
}

func dockerLogin() {
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if creds == "" {
		return
	}
	credsJson, err := ioutil.ReadFile(creds)
	if err != nil {
		log.Fatalf("There were errors reading the json key. %s", err.Error())
	}
	cmd := exec.Command("docker", "login", "-e", "1234@5678.com", "-u", "_json_key", "-p", string(credsJson), "https://gcr.io")
	fmt.Printf("Running docker login %s\n", cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env := os.Environ()
	cmd.Env = env
	err = cmd.Start()
	if err != nil {
		log.Fatalf("There were errors logging into docker. %s", err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("There were errors logging into docker. %s", err.Error())
	}
}

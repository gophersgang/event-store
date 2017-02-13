package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	//TODO generate the stub client???????
	pathToMainGo             = pathToServerDir + "/main.go"
	stubServerHeaderTemplate = `
package main
import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/vendasta/gosdks/logging"
	grpc_context "golang.org/x/net/context"
	"%s/api"
)
`
	stubServerDefinition   = "\ntype %s struct{}\n"
	stubFunctionDefinition = `
func (s *%s) %s(ctx grpc_context.Context, req *%s) (*%s, error) {
	//TODO: CHANGE THIS TO DO STUFF
	logging.Debugf(ctx, "You should change this")
	res := &%s{}
	return res, nil
}

`
)

//GenerateStubImplementations reads protobuf files from ./pb/*.pb.go and will create a files
//that contains empty server implementation for each protobuf file. It will also register the
//generated services in ./server/main.go
func GenerateStubImplementations(config MicroserviceConfig, protoFiles []string) {
	fmt.Printf("Generating Stub files in %s and augmenting %s\n", pathToServerDir, pathToMainGo)

	protoBufs := readProtoBufFiles()
	var services []Service
	for _, protoBuf := range protoBufs {
		createStubServerImplementation(config.GoPackageName, protoBuf)
		services = append(services, buildServicesFromProtoBuf(protoBuf)...)
	}
	augmentMainGoToRegisterServices(services, config.GoPackageName)

	fmt.Print("\n\nYou should now add your domain implementation to the files generated in pkg/\n\n")

	defer formatGeneratedCode()
}

type ProtoBuf struct {
	FileName string
	Content  []byte
}

func readProtoBufFiles() []ProtoBuf {
	var protoBufs []ProtoBuf
	if fileNames, err := filepath.Glob("./pb/*.pb.go"); err == nil {
		for _, fileName := range fileNames {
			b := justRead(fileName)
			protoBufs = append(protoBufs, ProtoBuf{FileName: fileName, Content: b})
		}

	} else {
		log.Fatalf("Error finding protobufs: %s", err.Error())
	}
	return protoBufs
}

func buildServicesFromProtoBuf(protoBuf ProtoBuf) []Service {
	var services []Service
	r := regexp.MustCompile("type (.*)Server interface {")
	if matches := r.FindAllSubmatch(protoBuf.Content, -1); matches != nil {
		for _, match := range matches {
			services = append(services, Service{Name: string(match[1:][0])})
		}
	}
	return services
}

func formatGeneratedCode() {
	format := exec.Command("goreturns", "-w", "server/")
	if err := format.Run(); err != nil {
		fmt.Printf("\n Error formating: %v\n", err)
	}
}

type Service struct {
	Name string
}

func (s *Service) getRegistrationString() string {
	return fmt.Sprintf("api.Register%sServer(grpcServer, &%sServer{})", s.Name, s.Name)
}

type StubServerData struct {
	Filename    string
	Servers     []StubServer
	PackageName string
}

func (ssd *StubServerData) generateFileContent() string {
	fileContent := fmt.Sprintf(stubServerHeaderTemplate, ssd.PackageName)
	for _, server := range ssd.Servers {
		fileContent = fileContent + server.generateFileContent()
	}
	return fileContent
}

type StubServer struct {
	Name      string
	Functions []StubFunction
}

func (ss *StubServer) generateFileContent() string {
	fileContent := fmt.Sprintf(stubServerDefinition, ss.Name)

	for _, function := range ss.Functions {
		fileContent = fileContent + function.generateFileContent(ss.Name)
	}
	return fileContent
}

type StubFunction struct {
	Name     string
	Request  string
	Response string
}

func (sf *StubFunction) generateFileContent(serverName string) string {
	return fmt.Sprintf(stubFunctionDefinition, serverName, sf.Name, sf.getFullRequestString(), sf.getFullResponseString(), sf.getFullResponseString())
}

func (sf *StubFunction) getFullRequestString() string {
	fullRequestString := "pb." + sf.Request
	if strings.HasPrefix(sf.Request, "google_protobuf") {
		fullRequestString = sf.Request
	}
	return fullRequestString
}

func (sf *StubFunction) getFullResponseString() string {
	fullResponseString := "pb." + sf.Response
	if strings.HasPrefix(sf.Response, "google_protobuf") {
		fullResponseString = sf.Response
	}
	return fullResponseString
}

func createStubServerImplementation(packageName string, protoBuf ProtoBuf) {
	servers := buildStubServersFromProtoBuf(protoBuf.Content)
	filename := getFileNameFromPath(protoBuf.FileName)
	stubServerData := StubServerData{Filename: filename, Servers: servers, PackageName: packageName}

	destination := fmt.Sprintf("%s/%s.go", pathToAPIDir, stubServerData.Filename)
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		if err := ioutil.WriteFile(destination, []byte(stubServerData.generateFileContent()), 0644); err != nil {
			log.Fatalf("Error writing %s/%s.go: %v", pathToAPIDir, stubServerData.Filename, err)
		}
	} else {
		fmt.Printf("\033[0;31m %s already exists, not generating the stub\033[0m \n", destination)
	}
}

func buildStubServersFromProtoBuf(b []byte) []StubServer {
	findServerDefinition := regexp.MustCompile(`type (.*Server) interface {[^}]+`)
	findingServerDifinitions := findServerDefinition.FindAllSubmatch(b, -1)
	var servers []StubServer
	for _, foundServerDefinitions := range findingServerDifinitions {
		fullServerDeclaration := foundServerDefinitions[0]
		functions := buildStubFunctionsFromServerDefinition(fullServerDeclaration)

		serverName := foundServerDefinitions[1]
		server := StubServer{Name: string(serverName), Functions: functions}
		servers = append(servers, server)
	}
	return servers
}

func buildStubFunctionsFromServerDefinition(serverDefinition []byte) []StubFunction {
	extractFunctionDefinitions := regexp.MustCompile(`(.+)\(context.Context, \*([^)]+)\) \(\*([^,]+), error\)`)
	findingFunctionDefinitions := extractFunctionDefinitions.FindAllSubmatch(serverDefinition, -1)
	var functions []StubFunction
	for _, foundFunctionDefinitions := range findingFunctionDefinitions {
		name := foundFunctionDefinitions[1]
		request := foundFunctionDefinitions[2]
		response := foundFunctionDefinitions[3]
		functions = append(functions, StubFunction{Name: string(name), Request: string(request), Response: string(response)})
	}
	return functions
}

func getFileNameFromPath(path string) string {
	findFileName := regexp.MustCompile(`pb/(.*).pb.go`)
	fileNameMatches := findFileName.FindAllStringSubmatch(path, -1)
	return fileNameMatches[0][1]
}

func insertRegisteringOfServers(mainGoText string, services []Service, packageName string) string {
	fmt.Println("Registering Servers")
	registerString := ""
	for _, service := range services {
		if !strings.Contains(mainGoText, service.getRegistrationString()) {
			registerString = registerString + service.getRegistrationString() + "\n"
		} else {
			fmt.Printf("Registration for %s already exists, skipping", service.Name)
		}
	}
	return strings.Replace(mainGoText, gRPCServerRegisterTag, registerString+gRPCServerRegisterTag, 1)
}

func insertProtobufImport(mainGoText string, packageName string) string {
	if !strings.Contains(mainGoText, packageName) {
		packageImport := fmt.Sprintf("import ( \n \"%s/pb\"", packageName)
		mainGoText = strings.Replace(mainGoText, "import (", packageImport, 1)
	} else {
		fmt.Printf("main.go already imports %s, skipping", packageName)
	}
	return mainGoText
}

func augmentMainGoToRegisterServices(services []Service, packageName string) {
	b := justRead(pathToMainGo)
	augmentedMainGo := insertRegisteringOfServers(string(b[:]), services, packageName)
	augmentedMainGo = insertProtobufImport(augmentedMainGo, packageName)

	if err := ioutil.WriteFile(pathToMainGo, []byte(augmentedMainGo), 0644); err != nil {
		log.Fatalf("Error writing %s", pathToMainGo)
	}
}

func justRead(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading %s", path)
	}
	return b
}

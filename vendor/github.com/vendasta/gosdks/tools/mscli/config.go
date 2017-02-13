package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"strings"

	"regexp"

	"github.com/ghodss/yaml"
)

type Network struct {
	GRPCHost       string `json:"grpcHost" yaml:"grpcHost"`
	GRPCPort       int `json:"grpcPort" yaml:"grpcPort"`
	HTTPSHost      string `json:"httpsHost" yaml:"httpsHost"`
	HTTPSPort      int    `json:"httpsPort" yaml:"httpsPort"`
	LoadBalancerIP string `json:"loadBalancerIp" yaml:"loadBalancerIp"`
}

type Scaling struct {
	MaxReplicas int32 `json:"maxReplicas" yaml:"maxReplicas"`
	MinReplicas int32 `json:"minReplicas" yaml:"minReplicas"`
	TargetCPU   int32 `json:"targetCPU" yaml:"targetCPU"`
}
type Resources struct {
	MemoryRequest string `json:"memoryRequest" yaml:"memoryRequest"`
	MemoryLimit   string `json:"memoryLimit" yaml:"memoryLimit"`
	CpuRequest    string `json:"cpuRequest" yaml:"cpuRequest"`
	CpuLimit      string `json:"cpuLimit" yaml:"cpuLimit"`
}

type AppConfig struct {
	EndpointsVersion string `json:"endpointsVersion" yaml:"endpointsVersion"`
}

type PodConfig struct {
	Secrets []Secret `json:"secrets" yaml:"secrets"`
	PodEnv  []Env    `json:"podEnv" yaml:"podEnv"`
}

type EnvironmentConfig struct {
	Name         string `json:"name" yaml:"name"`
	K8sContext   string `json:"k8sContext" yaml:"k8sContext"`
	K8sNamespace string `json:"k8sNamespace" yaml:"k8sNamespace"`

	Network `json:"network" yaml:"network"`
	Scaling `json:"scaling" yaml:"scaling"`
	Resources `json:"resources" yaml:"resources"`
	AppConfig `json:"appConfig" yaml:"appConfig"`
	PodConfig `json:"podConfig" yaml:"podConfig"`
}

type Env struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type Secret struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
}

type Redis struct {
	Password string `json:"password" yaml:"password"`
}

type Apps struct {
	Redis *Redis `json:"redis" yaml:"redis"`
}

type MicroserviceConfig struct {
	Name          string              `json:"name" yaml:"name"`
	GoPackageName string              `json:"goPackageName" yaml:"goPackageName"`
	Environments  []EnvironmentConfig `json:"environments" yaml:"environments"`
	Apps          Apps `json:"apps" yaml:"apps"`

	Debug         bool   `json:"-" yaml:"-"`
	Version       string `json:"-" yaml:"-"`
	Environment   string `json:"-" yaml:"-"`
}

type MicroserviceFile struct {
	Microservice MicroserviceConfig `json:"microservice" yaml:"microservice"`
}

func FailIfFileExists(fileName string) {
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		log.Fatalf("Error, %s exists.  Please remove and try again.", fileName)
	}
}

func EnsureDirExists(dirName string) {
	var err error
	var stat os.FileInfo

	if stat, err = os.Stat(dirName); err == nil {
		if stat.IsDir() {
			fmt.Printf("%s directory already exists, skipping creation\n", dirName)
			return
		}
		log.Fatalf("File %s exists, and is not a directory", dirName)
	} else {
		if os.IsNotExist(err) {
			if err = os.Mkdir(dirName, 0755); err != nil {
				log.Fatalf("Error creating directory %s: %s", dirName, err.Error())
			} else {
				fmt.Printf("%s directory created\n", dirName)
				return
			}
		}
	}
	log.Fatalf("Error encountered with creating directory %s: %s", dirName, err.Error())

}

func WriteConfigFile(config MicroserviceFile, name string) {
	var err error
	var data []byte

	//Serialize the config
	if data, err = yaml.Marshal(config); err != nil {
		log.Fatalf("Error marshalling YAML: %s", err.Error())
	}

	//Write it to file
	if err = ioutil.WriteFile(name, data, 0666); err != nil {
		log.Fatalf("Error writing file: %s", err.Error())
	}
}

func NewConfigFile() MicroserviceFile {
	return MicroserviceFile{
		Microservice: NewConfig(),
	}
}

// creates a new, default configuration file
func NewConfig() MicroserviceConfig {
	config := MicroserviceConfig{
		Environments: []EnvironmentConfig{
			NewEnvironment("test"),
			NewEnvironment("local"),
		},
	}

	config.GoPackageName = inferPackageNameFromDirectoryStructure()
	config.Name = inferNameFromPackageName(config.GoPackageName)
	return config
}

func inferPackageNameFromDirectoryStructure() string {
	defaultPackageName := "github.com/vendasta/my-project/my-service"
	r := regexp.MustCompile("github.com/vendasta.+")
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		if match := r.FindString(dir); match != "" {
			defaultPackageName = match
		}

	}
	return defaultPackageName
}

func inferNameFromPackageName(packageName string) string {
	splits := strings.Split(packageName, "/")
	possibleName := splits[len(splits) - 1]
	if possibleName == "" {
		possibleName = "MyService"
	}
	return possibleName
}

func NewEnvironment(name string) EnvironmentConfig {
	var k8sContext string
	var k8sNamespace string
	if name == "local" {
		k8sContext = "minikube"
		k8sNamespace = "default"
	} else {
		k8sContext = fmt.Sprintf("gke_repcore-prod_us-central1-c_vendasta-central")
		k8sNamespace = "sandbox"
	}
	var podEnvs []Env
	podEnvs = append(podEnvs, Env{Key: "ENV_VARIABLE", Value: "example"})

	var secrets []Secret
	// secrets = append(secrets, Secret{Name: "example-secret", MountPath: "/dev/null"}) //This example actually break the deployment

	return EnvironmentConfig{
		Name:          name,
		K8sNamespace:  k8sNamespace,
		K8sContext:    k8sContext,

		Resources: Resources{
			CpuLimit:      "250m",
			CpuRequest:    "100m",
			MemoryLimit:   "128Mi",
			MemoryRequest: "64Mi",
		},
		Network: Network{
			GRPCHost:      fmt.Sprintf("my-service-api-%s.vendasta-internal.com", name),
			HTTPSHost:      fmt.Sprintf("my-service-%s.vendasta-internal.com", name),
			GRPCPort:      11000,
			HTTPSPort:     11001,
		},
		PodConfig: PodConfig{
			PodEnv:        podEnvs,
			Secrets:       secrets,
		},
		Scaling: Scaling{
			MaxReplicas: 3,
			MinReplicas: 1,
			TargetCPU: 50,
		},
	}
}

func parseConfig(fileName string) MicroserviceConfig {
	var err error
	var data []byte
	var msf MicroserviceFile

	if data, err = ioutil.ReadFile(fileName); err != nil {
		log.Fatalf("Error could not open file: %s", err.Error())
	}

	if err = yaml.Unmarshal(data, &msf); err != nil {
		log.Fatalf("Error could not parse config: %s", err.Error())
	}
	ms := msf.Microservice

	ms.Name = strings.ToLower(ms.Name)

	//Apply the defaults
	for i := range ms.Environments {
		if ms.Environments[i].Network.GRPCPort == 0 {
			ms.Environments[i].Network.GRPCPort = 11000
		}
		if ms.Environments[i].Network.HTTPSPort == 0 {
			ms.Environments[i].Network.HTTPSPort = 11001
		}
	}

	return ms
}

func (ms *MicroserviceConfig) GetEnvironment() EnvironmentConfig {
	for i := range ms.Environments {
		if ms.Environments[i].Name == ms.Environment {
			return ms.Environments[i]
		}
	}
	log.Fatalf("Could not find environment %s", ms.Environment)
	return EnvironmentConfig{}
}

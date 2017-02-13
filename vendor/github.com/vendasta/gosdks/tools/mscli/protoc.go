package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildFromProtos(config MicroserviceConfig, protoSourceDir string) {
	validateProtoDir(protoSourceDir)

	var err error
	var wd string
	if wd, err = os.Getwd(); err != nil {
		log.Fatalf("Error getting work directory: %s", err.Error())
	}

	var files, paths []string
	if paths, err = filepath.Glob(fmt.Sprintf("%s/*.proto", protoSourceDir)); err != nil {
		log.Fatalf("Error expanding glob: %s", err.Error())
	}
	for i := range paths {
		r, _ := filepath.Rel(protoSourceDir, paths[i])
		files = append(files, r)
	}

	if config.Debug {
		fmt.Printf("Found the following proto files to be compiled: %q\n", files)
	}

	command := append([]string{"run", "--rm",
		"-v", fmt.Sprintf("%s:/src", protoSourceDir),
		"-v", fmt.Sprintf("%s/pb:/dest", wd),
		"vendasta/protoc-go",
		"--go_out=plugins=grpc,import_path=pb:/dest"}, files...)
	cmd := exec.CommandContext(context.Background(), "docker", command...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatalf("Error starting protoc: %s", err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("There were errors compiling protos: %s", err.Error())
	}

	GenerateStubImplementations(config, files)

}

func validateProtoDir(protoSourceDir string) {
	info, err := os.Stat(protoSourceDir)
	if err != nil {
		log.Fatalf("Error validating the proto path: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("The path to the protos must be a directory, not a file")
	}
}

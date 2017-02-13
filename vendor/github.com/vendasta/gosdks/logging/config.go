package logging

import (
	"errors"
	"cloud.google.com/go/logging"
	gce_metadata "cloud.google.com/go/compute/metadata"
	"google.golang.org/api/option"
	"context"
	"sync"
	"fmt"
	"strings"
)

var config *Config
var mut sync.Mutex

type Config struct {
	ProjectId     string
	Namespace     string
	PodName       string
	AppName       string
	TracingClient *Client
}

func (c *Config) BuildHeader(name string) string {
	return fmt.Sprintf("x-%s-%s", strings.ToLower(c.AppName), name)
}

func Initialize(gkeNamespace, podName, appName string) error {
	if gkeNamespace == "" || podName == "" || appName == "" {
		return errors.New("gkeNamespace, podName and appName must be supplied.")
	}
	mut.Lock()
	defer mut.Unlock()

	if config != nil {
		return errors.New("Logger has already been initialized.")
	}
	projectId := appName + "-local"
	if gce_metadata.OnGCE() {
		var err error
		projectId, err = gce_metadata.ProjectID(); if err != nil {
			return err
		}
	}

	tracingClient, err := NewClient(context.Background(), projectId); if err != nil {
		return err
	}

	config = &Config{
		Namespace: gkeNamespace,
		PodName: podName,
		AppName: appName,
		ProjectId: projectId,
		TracingClient: tracingClient,
	}
	ctx := context.Background()
	if gce_metadata.OnGCE() {
		client, err := logging.NewClient(ctx, projectId, option.WithGRPCConnectionPool(5)); if err != nil {
			return err
		}
		logger, err = newGkeLogger(config, client); if err != nil {
			return err
		}
	} else {
		logger, _ = newStdOutLogger(config)
	}
	return nil
}

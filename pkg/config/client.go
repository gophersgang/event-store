package config

import (
	"github.com/vendasta/gosdks/vstore"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"os"
)

var (
	// Namespace is the vstore namespace for this instance
	Namespace = vstore.Namespace(*vstore.Env(), "event-store")

	// OpportunityElasticIndex is used to specify the index when querying
	EventStoreElasticIndex string

	// ElasticIndexID is used when specifying the index via vstore
	ElasticIndexID = "elastic"
)

// IsLocal returns true if this instance is running locally
func IsLocal() bool {
	e := vstore.Env()
	return *e == vstore.Local || *e == vstore.Internal
}

// cache email
var email = ""

// GetCurrentLocalUserName returns the local user's name based on their google oauth email
func GetCurrentLocalUserName() string {
	if email != "" {
		return email
	}
	client, err := google.DefaultClient(context.Background(), oauth2.UserinfoEmailScope)
	service, err := oauth2.New(client)
	if err != nil {
		panic(err)
	}
	us := oauth2.NewUserinfoService(service)
	ui, err := us.Get().Do()
	if err != nil {
		panic(err)
	}
	email = ui.Email
	return email
}

// GetEnv returns the value of the environment variable specified by key, or fallback if no value is found
func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetServiceAccount returns the service account to use for this instance
func GetServiceAccount() string {
	if IsLocal() {
		return GetCurrentLocalUserName()
	}
	return GetEnv("SERVICE_ACCOUNT", "sales-opportunities-local@repcore-prod.iam.gserviceaccount.com")
}

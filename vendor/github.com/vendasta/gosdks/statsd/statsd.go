package statsd

import (
	"github.com/DataDog/datadog-go/statsd"
	"os"
	"fmt"
	"time"
	"errors"
	"github.com/vendasta/gosdks/config"
)

var client statsdInterface
var clientNotInitialized = errors.New("StatsD client has not initialized.")

func Initialize(statsNamespace string) error {
	ddAgentAddr := os.Getenv("DD_AGENT_ADDR")
	if ddAgentAddr == "" {
		ddAgentAddr = "dd-agent.default.svc.cluster.local:8125"
	}

	//use a fake client on local
	if config.IsLocal() {
		client = &fakeStatsD{}
		return nil
	}

	//use the datadog client on real environments
	c, err := statsd.New(ddAgentAddr); if err != nil {
		fmt.Printf("Error initializing statsd client. %s", err.Error())
		return err
	}
	client = &dataDogStatsD{c}
	client.SetNamespace(statsNamespace)
	return nil
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Gauge(name, value, tags, rate)
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Count(name, value, tags, rate)
}

// Histogram tracks the statistical distribution of a set of values.
func Histogram(name string, value float64, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Histogram(name, value, tags, rate)
}

// Decr is just Count of 1
func Decr(name string, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Decr(name, tags, rate)
}

// Incr is just Count of 1
func Incr(name string, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Incr(name, tags, rate)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Set(name, value, tags, rate)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.Timing(name, value, tags, rate)
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	if client == nil {
		return clientNotInitialized
	}
	return client.TimeInMilliseconds(name, value, tags, rate)
}

type statsdInterface interface {
	// Gauge measures the value of a metric at a particular time.
	Gauge(name string, value float64, tags []string, rate float64) error

	// Count tracks how many times something happened per second.
	Count(name string, value int64, tags []string, rate float64) error

	// Histogram tracks the statistical distribution of a set of values.
	Histogram(name string, value float64, tags []string, rate float64) error

	// Decr is just Count of 1
	Decr(name string, tags []string, rate float64) error

	// Incr is just Count of 1
	Incr(name string, tags []string, rate float64) error

	// Set counts the number of unique elements in a group.
	Set(name string, value string, tags []string, rate float64) error

	// Timing sends timing information, it is an alias for TimeInMilliseconds
	Timing(name string, value time.Duration, tags []string, rate float64) error

	// TimeInMilliseconds sends timing information in milliseconds.
	// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
	TimeInMilliseconds(name string, value float64, tags []string, rate float64) error

	// SetNamespace configures the prefix for all stats being pushed by this application.
	SetNamespace(string)
}

type dataDogStatsD struct {
	*statsd.Client
}

func (d *dataDogStatsD) SetNamespace(namespace string) {
	d.Client.Namespace = fmt.Sprintf("%s.", namespace)
}

package logging

import (
	"strings"
	"fmt"
	"time"
	"github.com/mattheath/kala/bigflake"
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"runtime"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	gce_metadata  "cloud.google.com/go/compute/metadata"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"strconv"
	"os"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/vendasta/gosdks/statsd"
)

type gkeLogger struct {
	logger *logging.Logger
	config *Config
	flake  *bigflake.Bigflake
}

func (l *gkeLogger) request(ctx context.Context, rl *RequestLog) {
	l.logRequest(ctx, rl)
}

func (l *gkeLogger) Debugf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Debug, f, a...)
}

func (l *gkeLogger) Infof(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Info, f, a...)
}

func (l *gkeLogger) Noticef(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Notice, f, a...)
}

func (l *gkeLogger) Warningf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Warning, f, a...)
}

func (l *gkeLogger) Errorf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Error, f, a...)
}

func (l *gkeLogger) Criticalf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Critical, f, a...)
}

func (l *gkeLogger) Alertf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Alert, f, a...)
}

func (l *gkeLogger) Emergencyf(ctx context.Context, f string, a ...interface{}) {
	l.log(ctx, logging.Emergency, f, a...)
}

func (l *gkeLogger) log(ctx context.Context, severity logging.Severity, f string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		filePieces := strings.Split(file, "/")
		file = strings.Join(filePieces[len(filePieces) - 2:], "/")
	}
	l.logger.Log(logging.Entry{
		Labels: GetMetadata(ctx),
		Timestamp: time.Now().UTC(),
		Severity: severity,
		Payload: fmt.Sprintf(f, a...),
		Operation: &logpb.LogEntryOperation{
			Id: GetRequestId(l.config, ctx),
			Producer: fmt.Sprintf("[%s:%d] ", file, line),
		},
	})
}

func (l *gkeLogger) logRequest(ctx context.Context, rl *RequestLog) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		filePieces := strings.Split(file, "/")
		file = strings.Join(filePieces[len(filePieces) - 2:], "/")
	}
	sev := logging.Info
	if rl.status >= 500 {
		sev = logging.Error
	}
	requestId := GetRequestId(l.config, ctx)
	l.logger.Log(logging.Entry{
		Labels: GetMetadata(ctx),
		Timestamp: time.Now().UTC(),
		Severity: sev,
		Operation: &logpb.LogEntryOperation{
			Id: requestId,
			Producer: fmt.Sprintf("[%s:%d] ", file, line),
		},
		Payload: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"@type": {&structpb.Value_StringValue{"type.googleapis.com/google.appengine.logging.v1.RequestLog"}},
				"traceId": {&structpb.Value_StringValue{rl.traceId}},
				"latency": {&structpb.Value_StringValue{fmt.Sprintf("%0.04fs", rl.latency.Seconds())}},
				"requestId": {&structpb.Value_StringValue{requestId}},
				"resource": {&structpb.Value_StringValue{rl.path}},
				"method": {&structpb.Value_StringValue{"POST"}},
				"httpVersion": {&structpb.Value_StringValue{"http/2"}},
				"status": {&structpb.Value_NumberValue{float64(rl.status)}},
				"responseSize": {&structpb.Value_StringValue{fmt.Sprintf("%d", rl.responseSize)}},
				"requestSize": {&structpb.Value_StringValue{fmt.Sprintf("%d", rl.requestSize)}},
				"userAgent": {&structpb.Value_StringValue{rl.userAgent}},
			},
		},
	})
	tags := []string{
		fmt.Sprintf("status:%d", rl.status),
		fmt.Sprintf("namespace:%s", l.config.Namespace),
		fmt.Sprintf("path:%s", rl.path),
	}
	go func() {
		statsd.Incr("gRPC", tags, 1)
		statsd.Histogram("gRPC.Latency", float64(rl.latency.Nanoseconds() / 1e6), tags, 1)
		statsd.Histogram("gRPC.RequestSize", float64(rl.requestSize), tags, 1)
		statsd.Histogram("gRPC.ResponseSize", float64(rl.responseSize), tags, 1)
	}()
}

func (l *gkeLogger) RequestId() string {
	for {
		f, err := l.flake.Mint();
		if err == bigflake.ErrSequenceOverflow {
			time.Sleep(time.Millisecond)
			continue
		} else if err != nil {
			Errorf(context.Background(), "Unable to use flake to generate a unique request id, returning empty string. Error: %s", err.Error())
			return ""
		}
		return f.String()
	}
}

func newGkeLogger(config *Config, client *logging.Client) (*gkeLogger, error) {
	projectId, err := gce_metadata.ProjectID(); if err != nil {
		return nil, err
	}
	zone, err := gce_metadata.Zone(); if err != nil {
		return nil, err
	}
	clusterName, err := gce_metadata.InstanceAttributeValue("cluster-name"); if err != nil {
		return nil, err
	}
	instanceId, err := gce_metadata.InstanceID(); if err != nil {
		return nil, err
	}
	instanceIdInt, err := strconv.Atoi(instanceId); if err != nil {
		return nil, err
	}
	labels := map[string]string{
		"project_id": projectId,
		"cluster_name": clusterName,
		"namespace_id": config.Namespace,
		"instance_id": instanceId,
		"pod_id": os.Getenv("HOSTNAME"),
		"container_name": config.AppName,
		"zone": zone,
	}
	mr := &mrpb.MonitoredResource{Type: "container", Labels: labels}
	workerId := uint64(instanceIdInt) & uint64((1 << 31) - 1); //take 32 bits from the 64-bit integer.
	flake, _ := bigflake.New(workerId)
	client.OnError = func(err error) {
		fmt.Printf("Error flushing logs: %s", err.Error())
	}
	logger := gkeLogger{
		config: config,
		flake: flake,
		logger: client.Logger(config.AppName, logging.CommonResource(mr)),
	}
	return &logger, nil
}

type RequestLog struct {
	path         string
	method       string
	requestSize  int64
	responseSize int64
	requestId    string
	status       int32
	latency      time.Duration
	traceId      string
	labels       map[string]string
	userAgent    string
}

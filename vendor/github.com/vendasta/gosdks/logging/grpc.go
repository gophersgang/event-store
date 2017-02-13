package logging

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"net/http"
	"time"
	"fmt"
	"github.com/golang/protobuf/proto"
	gce_metadata  "cloud.google.com/go/compute/metadata"
	"io"
)

func Interceptor() grpc.UnaryServerInterceptor {
	if gce_metadata.OnGCE() {
		i := &grpcInterceptor{config: config, logger: GetLogger()}
		return i.UnaryServerInterceptor
	}
	return PassThroughUnaryServerInterceptor
}

func ClientInterceptor() (grpc.UnaryClientInterceptor) {
	if gce_metadata.OnGCE() {
		i := &grpcInterceptor{config: config, logger: GetLogger()}
		return i.UnaryClientInterceptor
	}
	return PassThroughUnaryClientInterceptor
}

func ClientStreamInterceptor() (grpc.StreamClientInterceptor) {
	if gce_metadata.OnGCE() {
		i := &grpcInterceptor{config: config, logger: GetLogger()}
		return i.StreamClientInterceptor
	}
	return PassThroughStreamClientInterceptor
}

type grpcInterceptor struct {
	config *Config
	logger Logger
}

func PassThroughUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return handler(ctx, req)
}

func PassThroughUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, reply, cc, opts...)
}

func PassThroughStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return streamer(ctx, desc, cc, method, opts...)
}

func (g *grpcInterceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now().UTC()
	ctx = SetValue(ctx, g.config.BuildHeader("request-id"), g.logger.RequestId())
	ctx = SetTimeValue(ctx, g.config.BuildHeader("start-time"), start)

	span := g.config.TracingClient.SpanFromContext(ctx, info.FullMethod)
	defer span.Finish()
	ctx = NewContext(ctx, span)

	resp, err := handler(ctx, req)

	end := time.Now().UTC()
	ctx = SetTimeValue(ctx, g.config.BuildHeader("end-time"), end)
	ctx = SetValue(ctx, g.config.BuildHeader("took"), fmt.Sprintf("%dms", (end.Sub(start).Nanoseconds() / 1e6)))
	statusCode := HTTPStatusFromCode(grpc.Code(err))
	if statusCode >= 500 {
		Errorf(ctx, "Error serving request with code %d Error: %s", grpc.Code(err), err.Error())
	}
	traceID := ""
	if span != nil {
		traceID = span.TraceID()
	}
	responseSize := 0
	p, ok := resp.(proto.Message); if ok {
		responseSize = proto.Size(p)
	}
	contentLength, _ := GetIntValue(ctx, "content-length")
	userAgent, _ := GetFirstValue(ctx, "user-agent")
	logRequest(ctx, &RequestLog{
		path: info.FullMethod,
		method: "POST",
		requestId: GetRequestId(g.config, ctx),
		responseSize: int64(responseSize),
		status: statusCode,
		latency: end.Sub(start),
		traceId: traceID,
		requestSize: contentLength,
		labels: GetMetadata(ctx),
		userAgent: userAgent,
	})
	return resp, err
}

func (g *grpcInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	span := FromContext(ctx)
	childSpan := span.NewRemoteChild(ctx, method)
	defer childSpan.Finish()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	Debugf(context.Background(), "Made call to %s. Took %d milliseconds.", method, end.Sub(start).Nanoseconds() / 1e6)
	return err
}

func (g *grpcInterceptor) StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	span := FromContext(ctx)
	childSpan := span.NewRemoteChild(ctx, method)

	cs, err := streamer(ctx, desc, cc, method, opts...); if err != nil {
		childSpan.Finish()
		return nil, err
	}

	return &monitoredClientStream{cs, childSpan}, nil
}

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct {
	grpc.ClientStream
	span *Span
}

func (s *monitoredClientStream) SendMsg(m interface{}) error {
	return s.ClientStream.SendMsg(m)
}

func (s *monitoredClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err == io.EOF {
		s.span.Finish()
	} else {
		s.span.Finish()
	}
	return err
}

func HTTPStatusFromCode(code codes.Code) int32 {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusForbidden
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Printf("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

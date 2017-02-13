package logging

import (
	"golang.org/x/net/context"
)

type Logger interface {
	request(ctx context.Context, r *RequestLog)
	Debugf(ctx context.Context, f string, a ...interface{})
	Infof(ctx context.Context, f string, a ...interface{})
	Warningf(ctx context.Context, f string, a ...interface{})
	Errorf(ctx context.Context, f string, a ...interface{})
	Criticalf(ctx context.Context, f string, a ...interface{})
	Alertf(ctx context.Context, f string, a ...interface{})
	Emergencyf(ctx context.Context, f string, a ...interface{})
	RequestId() string
}

var logger Logger = nil;

func GetLogger() Logger {
	if logger == nil {
		logger = &stdoutLogger{config: &Config{}}
	}
	return logger
}

func logRequest(ctx context.Context, r *RequestLog) {
	GetLogger().request(ctx, r)
}

func Debugf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Debugf(ctx, f, a...)
}

func Infof(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Infof(ctx, f, a...)
}

func Warningf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Warningf(ctx, f, a...)
}

func Errorf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Errorf(ctx, f, a...)
}

func Criticalf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Criticalf(ctx, f, a...)
}

func Alertf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Alertf(ctx, f, a...)
}

func Emergencyf(ctx context.Context, f string, a ...interface{}) {
	GetLogger().Emergencyf(ctx, f, a...)
}

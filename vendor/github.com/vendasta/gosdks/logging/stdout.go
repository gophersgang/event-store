package logging

import (
	"strings"
	"fmt"
	"golang.org/x/net/context"
	"math/rand"
	"runtime"
	"cloud.google.com/go/logging"
)

const (
	ColorRed = iota + 30
	ColorGreen
	ColorYellow
	ColorMagenta
	ColorCyan
)

type color int

var (
	colors = []string{
		logging.Debug: ColorSeq(ColorMagenta),
		logging.Info:    ColorSeq(ColorRed),
		logging.Warning:  ColorSeq(ColorYellow),
		logging.Error:   ColorSeq(ColorGreen),
		logging.Critical:    ColorSeq(ColorCyan),
		logging.Alert:    ColorSeq(ColorCyan),
		logging.Emergency:    ColorSeq(ColorCyan),
	}
)

func ColorSeq(color color) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

type stdoutLogger struct {
	config *Config
}

func (l *stdoutLogger) request(ctx context.Context, r *RequestLog) {
	l.Log(ctx, logging.Debug, "Served gRPC request for handler %s with code %d", r.path, r.status)
}
func (l *stdoutLogger) Debugf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Debug, f, a...)
}

func (l *stdoutLogger) Infof(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Info, f, a...)
}

func (l *stdoutLogger) Warningf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Warning, f, a...)
}

func (l *stdoutLogger) Errorf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Error, f, a...)
}

func (l *stdoutLogger) Criticalf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Critical, f, a...)
}

func (l *stdoutLogger) Alertf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Alert, f, a...)
}

func (l *stdoutLogger) Emergencyf(ctx context.Context, f string, a ...interface{}) {
	l.Log(ctx, logging.Emergency, f, a...)
}

func (l *stdoutLogger) Log(ctx context.Context, level logging.Severity, f string, a ...interface{}) {
	col := colors[level]
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		filePieces := strings.Split(file, "/")
		file = strings.Join(filePieces[len(filePieces) - 2:], "/")
	}
	if !strings.HasSuffix(f, "\n") {
		f = f + "\n"
	}
	md := GetMetadata(ctx)
	prefix := fmt.Sprintf("%s%-6s%50s:%-4d\033[0m", col, level.String(), file, line)
	metadata := ""
	for key, val := range md {
		metadata = metadata + fmt.Sprintf("[%s:%s]", key, val)
	}
	if metadata != "" {
		fmt.Printf(metadata + "\n")
	}
	fmt.Printf(prefix + " " + f, a...)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (l *stdoutLogger) RequestId() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newStdOutLogger(config *Config) (*stdoutLogger, error) {
	return &stdoutLogger{config: config}, nil
}

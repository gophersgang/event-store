package logging

import (
	"net/http"
	"fmt"
	"time"
)

func newLoggedResponse(w http.ResponseWriter) *loggedResponse {
	return &loggedResponse{w, 200, 0}
}

type loggedResponse struct {
	http.ResponseWriter
	status int
	length int
}

func (l *loggedResponse) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *loggedResponse) Write(b []byte) (n int, err error){
	n, err = l.ResponseWriter.Write(b)
	l.length += n
	return
}

func HTTPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		requestId := GetLogger().RequestId()
		requestUrl := r.URL.Host + "/" + r.URL.Path
		ctx = SetValue(ctx, config.BuildHeader("request-id"), requestId)
		ctx = SetValue(ctx, config.BuildHeader("request-url"), requestUrl)

		r = r.WithContext(ctx)

		Infof(r.Context(), fmt.Sprintf("Serving request %s . request-id %s", requestUrl, requestId))

		start := time.Now()

		lw := newLoggedResponse(w)

		h.ServeHTTP(lw, r)
		end := time.Now()

		logRequest(ctx, &RequestLog{
			path: r.URL.Path,
			method: r.Method,
			requestId: requestId,
			responseSize: int64(lw.length),
			status: int32(lw.status),
			latency: end.Sub(start),
			requestSize: r.ContentLength,
			labels: GetMetadata(ctx),
			userAgent: userAgent,
		})
	})
}

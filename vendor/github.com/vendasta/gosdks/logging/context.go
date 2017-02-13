package logging

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GetMetadata(ctx context.Context) map[string]string {
	md, _ := metadata.FromContext(ctx);
	m := map[string]string{}
	for k, v := range md {
		if len(v) == 0 {
			continue
		} else {
			m[k] = strings.Join(v, ",")
		}
	}
	return m
}

func GetRequestId(config *Config, ctx context.Context) string {
	request_id, _ := GetFirstValue(ctx, config.BuildHeader("request-id"))
	return request_id
}

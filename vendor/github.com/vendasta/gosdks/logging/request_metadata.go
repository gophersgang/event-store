package logging

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"time"
	"strconv"
)

func GetFirstValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromContext(ctx); if !ok {
		return "", false
	}
	val, ok := md[key]; if !ok {
		return "", false
	}
	if len(val) >= 1 {
		return val[0], true
	}
	return "", false
}

func SetValue(ctx context.Context, key, value string) context.Context {
	return mergeMetadata(ctx, metadata.Pairs(key, value))
}

func SetTimeValue(ctx context.Context, key string, val time.Time) context.Context {
	return mergeMetadata(ctx, metadata.Pairs(key, strconv.Itoa(int(val.UnixNano()))))
}

func GetIntValue(ctx context.Context, key string) (int64, error) {
	val, _ := GetFirstValue(ctx, key)
	if val == "" {
		return 0, nil
	}
	intVal, err := strconv.Atoi(val)
	return int64(intVal), err
}

// mergeMetadata returns a context populated by the existing metadata, if any,
// joined with internal metadata.
func mergeMetadata(ctx context.Context, md metadata.MD) context.Context {
	mdCopy, _ := metadata.FromContext(ctx)
	return metadata.NewContext(ctx, metadata.Join(mdCopy, md))
}

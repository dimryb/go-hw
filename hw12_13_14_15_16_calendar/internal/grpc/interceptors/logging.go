package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

func UnaryLoggerInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		log.Infof("gRPC request: %s", info.FullMethod)
		log.Debugf("Request payload: %+v", req)

		resp, err := handler(ctx, req)

		duration := time.Since(start).Milliseconds()
		status := "success"
		if err != nil {
			status = "error"
		}

		log.Infof("gRPC finished: method=%s duration=%dms status=%s error=%v",
			info.FullMethod, duration, status, err)

		return resp, err
	}
}

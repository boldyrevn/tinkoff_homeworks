package grpcserver

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

// LoggingServerInterceptor implements logging for GRPC server
func LoggingServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Request - Method:%s; Duration:%v; Error:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return resp, err
}

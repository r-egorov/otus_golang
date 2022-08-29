package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func logInterceptor(
	log server.Logger,
) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	interceptor := func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		var remoteAddr string
		if p, ok := peer.FromContext(ctx); ok {
			remoteAddr = p.Addr.String()
		}
		res, err := handler(ctx, req)
		log.Info(fmt.Sprintf(`%s [%s] %s %s %s gRPC-Call"`,
			remoteAddr,
			start.Format("01/Jan/2001:12:00:00 +0300"),
			info.FullMethod,
			req,
			time.Since(start)/time.Second,
		))
		return res, err
	}
	return interceptor
}

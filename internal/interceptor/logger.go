package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	// Get client IP if available
	p, _ := peer.FromContext(ctx)
	clientIP := "unknown"
	if p != nil {
		clientIP = p.Addr.String()
	}

	// Log the incoming request
	slog.Info("incoming request",
		"method", info.FullMethod,
		"client_ip", clientIP,
	)

	// Handle the request
	resp, err := handler(ctx, req)
	
	// Log the response
	duration := time.Since(start)
	if err != nil {
		slog.Error("request failed",
			"method", info.FullMethod,
			"client_ip", clientIP,
			"duration", duration,
			"error", err,
		)
	} else {
		slog.Info("request completed",
			"method", info.FullMethod,
			"client_ip", clientIP,
			"duration", duration,
		)
	}

	return resp, err
}

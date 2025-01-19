package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/malinatrash/kartinki-auth/internal/config"
	"github.com/malinatrash/kartinki-auth/internal/kafka"
	"github.com/malinatrash/kartinki-auth/internal/repository/postgres"
	"github.com/malinatrash/kartinki-auth/internal/service"
	pb "github.com/malinatrash/kartinki-proto/gen/go/auth_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.Default()
	cfg := config.Load()

	host := strings.TrimPrefix(cfg.Host, "http://")

	addr := fmt.Sprintf("%s:%s", host, cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen", "error", err)
		panic(err)
	}

	repo, err := postgres.NewRepository(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	)
	if err != nil {
		logger.Error("failed to create repository", "error", err)
		panic(err)
	}
	defer repo.Close()

	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, logger, repo)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.ReadUsers(ctx)

	go func() {
		<-ctx.Done()
		cancel()
	}()

	authService := service.NewAuthService(cfg.JWTSecret, repo, logger)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authService)

	reflection.Register(grpcServer)

	logger.Info("auth service started", "address", addr)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err)
		panic(err)
	}
}

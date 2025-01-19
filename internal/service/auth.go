package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malinatrash/kartinki-auth/internal/repository/postgres"
	pb "github.com/malinatrash/kartinki-proto/gen/go/auth_service/v1"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	jwtSecret string
	repo      *postgres.Repository
	logger    *slog.Logger
	UserDeleter
	UserGetter
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
}

type UserGetter interface {
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
}

func NewAuthService(jwtSecret string, repo *postgres.Repository, logger *slog.Logger) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		repo:      repo,
		logger:    logger,
	}
}

// GetUser authenticates a user by their secret and returns a JWT token
// containing the user's ID, username and avatar. The token is valid for 24 hours.
func (s *AuthService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.logger.Info("processing GetUser request", "secret", req.Secret)

	user, err := s.repo.GetUserBySecret(req.Secret)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	s.logger.Info("generated JWT token", "token", token)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("failed to generate JWT token", "error", err)
		return nil, err
	}

	s.logger.Info("user authenticated successfully", "user_id", user.ID)
	return &pb.GetUserResponse{
		Jwt: tokenString,
		User: &pb.User{
			Id:       fmt.Sprint(user.ID),
			Username: user.Username,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.Info("processing DeleteUser request", "secret", req.Secret)
	secret := req.Secret
	if ok, err := s.repo.DeleteUser(ctx, secret); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("user not found")
	}
	s.logger.Info("user deleted successfully")
	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

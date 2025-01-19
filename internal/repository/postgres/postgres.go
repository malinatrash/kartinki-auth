package postgres

import (
	"fmt"
	"log/slog"

	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserGetter interface {
	GetUser(ctx context.Context, secret string) (*User, error)
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, secret string) (bool, error)
}

type Repository struct {
	db *gorm.DB
	ug UserGetter
	dg UserDeleter
}

func NewRepository(host, port, user, password, dbname string) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run automigrations
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run automigrations: %w", err)
	}

	slog.Info("connected to database and ran migrations successfully")
	return &Repository{db: db}, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
	)
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

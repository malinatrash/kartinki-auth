package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/malinatrash/kartinki-auth/internal/repository/postgres"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	repo   *postgres.Repository
	logger *slog.Logger
}

func NewConsumer(brokers []string, topic, groupID string, logger *slog.Logger, repo *postgres.Repository) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		logger: logger,
		repo:   repo,
	}
}

func (c *Consumer) ReadUsers(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("Error reading message", "error", err)
			continue
		}

		c.logger.Info("Received Kafka message",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"key", string(msg.Key),
			"value", string(msg.Value))

		var user *User
		if err := json.Unmarshal(msg.Value, &user); err != nil {
			c.logger.Error("Error unmarshalling user", 
				"error", err,
				"value", string(msg.Value))
			continue
		}

		repoUser := postgres.User{
			ID:       uint(user.Id),
			Username: user.Username,
			Avatar:   user.Avatar,
			Secret:   user.Secret,
		}

		if err := c.repo.CreateUser(&repoUser); err != nil {
			c.logger.Error("Error creating user", 
				"error", err,
				"user", user)
			continue
		}

		c.logger.Info("User created successfully", 
			"id", user.Id,
			"username", user.Username,
			"key", string(msg.Key))
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}

package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Host      string `env:"AUTH_HOST"`
	Port      string `env:"AUTH_PORT"`
	JWTSecret string `env:"JWT_SECRET"`

	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`

	KafkaBrokers []string `env:"KAFKA_BROKERS"`
	KafkaTopic   string   `env:"KAFKA_TOPIC"`
	KafkaGroupID string   `env:"KAFKA_GROUP_ID"`
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := &Config{}
	cfg.Host = os.Getenv("AUTH_HOST")
	cfg.Port = os.Getenv("AUTH_PORT")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	cfg.PostgresHost = os.Getenv("POSTGRES_HOST")
	cfg.PostgresPort = os.Getenv("POSTGRES_PORT")
	cfg.PostgresUser = os.Getenv("POSTGRES_USER")
	cfg.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	cfg.PostgresDB = os.Getenv("POSTGRES_DB")

	cfg.KafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	cfg.KafkaTopic = os.Getenv("KAFKA_TOPIC")
	cfg.KafkaGroupID = os.Getenv("KAFKA_GROUP_ID")

	return cfg
}

package facility

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN"`

	TimeLimit int `env:"TIME_LIMIT" envDefault:"5"`

	S3Region     string `env:"S3_REGION" envDefault:"eu-north-1"`
	S3AccessKey  string `env:"S3_ACCESS_KEY"`
	S3SecretKey  string `env:"S3_SECRET_KEY"`
	S3BucketName string `env:"S3_BUCKET_NAME"`
}

// NewConfig creates a new Config
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readFromEnvironment(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// readFromEnvironment reads the settings from environment variables.
func (c *Config) readFromEnvironment() error {
	return env.Parse(c)
}

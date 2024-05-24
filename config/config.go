package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config represents the application configuration
type Config struct {
	DB         DBConfig `env:",prefix=DB_,required"`
	JWTSecret  string   `env:"JWT_SECRET"`
	BcryptSalt int      `env:"BCRYPT_SALT"`
	S3Region   string   `env:"AWS_REGION"`
	S3Bucket   string   `env:"AWS_S3_BUCKET_NAME"`
	S3AcessKey string   `env:"AWS_ACCESS_KEY_ID"`
	S3Secret   string   `env:"AWS_SECRET_ACCESS_KEY"`
}

type DBConfig struct {
	Name     string `env:"NAME"`
	Port     string `env:"PORT"`
	Host     string `env:"HOST"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Params   string `env:"PARAMS"`
}

func Load(ctx context.Context) (*Config, error) {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ConnectionURL returns the connection URL for the database
func (c DBConfig) ConnectionString() string {
	params := strings.ReplaceAll(c.Params, `"`, "")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Name, params)
}

package config

import (
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/caarlos0/env"
	"github.com/danielealbano/svdb/shared/collection"
	_ "github.com/joho/godotenv/autoload"
	"github.com/phuslu/log"
	"time"
)

type Config struct {
	Host                       string `env:"HOST" envDefault:"0.0.0.0"`
	Port                       int    `env:"PORT" envDefault:"3000"`
	LogLevel                   string `env:"LOG_LEVEL" envDefault:"info"`
	CollectionQuantization     string `env:"COLLECTION_QUANTIZATION" envDefault:"F32"`
	CollectionMetric           string `env:"COLLECTION_METRIC" envDefault:"Cosine"`
	CollectionVectorDimensions uint   `env:"COLLECTION_VECTOR_DIMENSIONS" envDefault:"128"`
	ShardPath                  string `env:"SHARD_PATH"`
	ShardWriteable             bool   `env:"SHARD_WRITEABLE" envDefault:"false"`
	ShardMaxSize               string `env:"SHARD_MAX_SIZE" envDefault:"1GB"`
	ShardAutoSync              bool   `env:"SHARD_AUTO_SYNC" envDefault:"false"`
	ShardAutoSyncInterval      string `env:"SHARD_AUTO_SYNC_INTERVAL" envDefault:"1m"`
}

func ParseShardMaxSize(size string) (uint, error) {
	v, err := datasize.ParseString(size)

	if err != nil {
		return 0, fmt.Errorf("failed to parse size: %w", err)
	}

	return uint(v), nil
}

func FromEnv() (*Config, error) {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func ValidateConfig(config *Config) error {
	var err error
	var maxSize uint
	var interval time.Duration

	if config.Host == "" {
		return fmt.Errorf("host is required")
	}
	if config.Port <= 0 {
		return fmt.Errorf("port must be greater than 0")
	}
	if config.ShardPath == "" {
		return fmt.Errorf("shard path is required")
	}

	if log.ParseLevel(config.LogLevel) > log.PanicLevel {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	if config.CollectionVectorDimensions <= 0 {
		return fmt.Errorf("collection vector dimensions must be greater than 0")
	}

	if _, err = shared_collection.ParseQuantization(config.CollectionQuantization); err != nil {
		return fmt.Errorf("invalid collection quantization: %s", config.CollectionQuantization)
	}

	if _, err = shared_collection.ParseMetric(config.CollectionMetric); err != nil {
		return fmt.Errorf("invalid collection metric: %s", config.CollectionMetric)
	}

	if config.ShardWriteable {
		maxSize, err = ParseShardMaxSize(config.ShardMaxSize)
		if err != nil {
			return fmt.Errorf("failed to parse the shard max size: %w", err)
		}

		if maxSize <= 0 {
			return fmt.Errorf("shard max size must be greater than 0")
		}

		if config.ShardAutoSync {
			interval, err = time.ParseDuration(config.ShardAutoSyncInterval)
			if err != nil {
				return fmt.Errorf("failed to parse the shard auto sync interval: %w", err)
			}

			if interval <= 0 {
				return fmt.Errorf("shard auto sync interval must be greater than 0")
			}
		}
	}

	return nil
}

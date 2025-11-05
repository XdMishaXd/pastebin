package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
}

type Kafka struct {
	Host  string `yaml:"host" env-default:"localhost"`
	Port  int    `yaml:"post" env-default:"9092"`
	Topic string `yaml:"topic" env-required:"true"`
}

type Hash struct {
	HashRate   int `yaml:"hash_rate" env-required:"true"`
	HashLength int `yaml:"hash_length" env-required:"true"`
	BatchSize  int `yaml:"batch_size" env-default:"1"`
}

func MustLoad() *Config {
	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return &cfg
}

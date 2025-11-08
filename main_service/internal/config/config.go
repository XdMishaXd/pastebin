package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	MySQL      `yaml:"mysql"`
	Kafka      `yaml:"kafka"`
	MinIO      `yaml:"minio"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type MinIO struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Bucket   string `yaml:"bucket" env-required:"true"`
}

type MySQL struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"dbname" env-required:"true"`
}

type Kafka struct {
	Addr    string `yaml:"addr" env-default:"kafka:9092"`
	Topic   string `yaml:"topic" env-required:"true"`
	GroupID string `yaml:"groupid" env-default:"1"`
}

func MustLoad(configPath string) *Config {
	// проверка существования файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return &cfg
}

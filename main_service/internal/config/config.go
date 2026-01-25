package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DefaultTTL int    `yaml:"default_ttl" env-default:"1"`
	HTTPServer `yaml:"http_server"`
	MySQL      `yaml:"mysql"`
	Kafka      `yaml:"kafka"`
	MinIO      `yaml:"minio"`
	Redis      `yaml:"redis"`
	Swagger    `yaml:"swagger"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Swagger struct {
	Username string `yaml:"username" env-default:"admin"`
	Password string `yaml:"password" env-default:"admin"`
	Enabled  bool   `yaml:"enabled" env-default:"false"`
}

type MinIO struct {
	Endpoint string `yaml:"endpoint" env-default:"localhost:9000"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Bucket   string `yaml:"bucket" env-required:"true"`
	UseSSL   bool   `yaml:"useSSL" env-default:"false"`
}

type MySQL struct {
	DSN string `yaml:"dsn" env-required:"true"`
}

type Kafka struct {
	Addr  string `yaml:"addr" env-default:"kafka:9092"`
	Topic string `yaml:"topic" env-required:"true"`
}

type Redis struct {
	Addr                string `yaml:"addr" env-default:"redis:6379"`
	Db                  int    `yaml:"db" env-default:"1"`
	PopularityThreshold int64  `yaml:"popularity_threshold" env-default:"500"`
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

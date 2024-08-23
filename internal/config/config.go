package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage            StorageConfig `yaml:"storage"`
	HTTPServer         HTTPConfig    `yaml:"http-server"`
	SMTPMailer         SMTPConfig    `yaml:"smtp-mailer"`
	ActivationTokenTTL time.Duration `yaml:"activation-token-ttl" env-default:"2h"`
	RabbitMQConnString string        `yaml:"rabbitmq-conn-string" env-required:"true"`
	OpenAIKEY          string        `yaml:"openai-key" env:"ARCUS_OPENAI_KEY" env-required:"true"`
	Env                string        `yaml:"env" env-default:"development"`
}

type HTTPConfig struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-default:"8080"`
}

type SMTPConfig struct {
	Host          string `yaml:"host" env-rquired:"true"`
	Port          int    `yaml:"port" env-default:"25"`
	Username      string `yaml:"username" env-rquired:"true"`
	Password      string `yaml:"password" env-rquired:"true"`
	SenderAddress string `yaml:"sender-address" env-required:"true"`
}

type StorageConfig struct {
	DSN string `yaml:"dsn" env-required:"true" env:"ARCUS_POSTGRES_DSN"`
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic("error loading config: " + err.Error())
	}
	return cfg
}

func Load() (*Config, error) {
	configPath, err := fetchConfigPath()
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fetchConfigPath() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config-path", os.Getenv("ARCUS_CONFIG_PATH"), "App configuration file path")
	flag.Parse()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("cannot find config path: %w", err)
	}
	return configPath, nil
}

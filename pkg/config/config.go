package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:env env-required:"local"`
	StoragePath string        `yaml:storage_path env-required:"true"`
	TokenTll    time.Duration `yaml:token_tll env-required:"true"`
	Grpc        GrpcConf      `yaml:grpc`
}

type GrpcConf struct {
	Port    int           `yaml:port`
	Timeout time.Duration `yaml:timeout`
}

func MustLoadConfig() *Config {
	path := fetchConfgPath()

	if path == "" {
		panic("config path is not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config file: " + err.Error())
	}
	return &cfg
}

func fetchConfgPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")

	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

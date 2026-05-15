package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string        `yaml:"env" env-default:"local"`
	TokenTTL  time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC      GRPC          `yaml:"grpc"`
	DataStore DataStore     `yaml:"datastore"`
}

type DataStore struct {
	Storage Storage `yaml:"storage"`
	Cache   Cache   `yaml:"cache"`
}

type Storage struct {
	Host     string `yaml:"host" env:"STORAGE_HOST"`
	Port     int    `yaml:"port" env:"STORAGE_PORT"`
	User     string `yaml:"user" env:"STORAGE_USER"`
	DBName   string `yaml:"dbname" env:"STORAGE_DBNAME"`
	Password string `yaml:"password" env:"STORAGE_PASSWORD"`
	SSLMode  string `yaml:"sslmode"`
}

type Cache struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	cfgPath := fetchConfigPath()

	if cfgPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config file does not exist: " + cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config_path", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

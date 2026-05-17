package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string        `yaml:"env" env:"ENV" env-default:"local"`
	TokenTTL  time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
	GRPC      GRPC          `yaml:"grpc"`
	DataStore DataStore     `yaml:"datastore"`
	Clients   ClientsConfig `yaml:"clients"`
}

type DataStore struct {
	Storage Storage `yaml:"storage"`
	Cache   Cache   `yaml:"cache"`
}

type Storage struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	DBName   string `yaml:"dbname" env:"DB_NAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
}

type Cache struct {
	Addr        string        `yaml:"addr" env:"REDIS_ADDR"`
	Password    string        `yaml:"password" env:"REDIS_PASSWORD"`
	DB          int           `yaml:"db" env:"REDIS_DB"`
	MaxRetries  int           `yaml:"max_retries" env:"REDIS_MAX_RETRIES"`
	DialTimeout time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT"`
	Timeout     time.Duration `yaml:"timeout" env:"REDIS_TIMEOUT"`
}

type GRPC struct {
	Port    int           `yaml:"port" env:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10h"`
}

type Broker struct {
	Address      string `yaml:"address" env:"KAFKA_ADDRESS"`
	TopicName    string `yaml:"topic" env:"KAFKA_TOPIC"`
	Network      string `yaml:"network" env:"KAFKA_NETWORK"`
	Partitions   int    `yaml:"partitions" env:"KAFKA_PARTITIONS"`
	Replications int    `yaml:"replications" env:"KAFKA_REPLICATIONS"`
	GroupID      string `yaml:"group_id" env:"KAFKA_GROUP_ID"`
}

type ClientsConfig struct {
	Broker Broker `yaml:"kafka"`
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

package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

// New Read environment variables and return a Config struct else defaults to localhost and port 8080
func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("config.env", &cfg)
	if err != nil {
		cfg = defaultConfig()
	}
	return &cfg
}
func defaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: "8080",
	}
}

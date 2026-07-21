package config

import "os"

type Config struct {
	Port      string
	StaticDir string
}

func Load() *Config {
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:      port,
		StaticDir: "../../../../client",
	}
}

package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/online_courses?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return &Config{
		DatabaseURL: dsn,
		Port:        port,
	}
}

// Package config
package config

import "time"

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JwtConfig
}

type AppConfig struct {
	Port      string
	Version   string
	RateLimit int
}

type DBConfig struct {
	URI      string
	Name     string
	Username string
	Password string
}

type JwtConfig struct {
	Secret     string
	Expiration time.Duration
}

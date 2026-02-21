package config

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultServerHost         = "localhost"
	DefaultServerPort         = 9000
	DefaultServerReadTimeout  = 30 * time.Second
	DefaultServerWriteTimeout = 30 * time.Second
	DefaultLogLevel           = "info"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN string
}

type LogConfig struct {
	Level string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	MFATempTTL     time.Duration
	Issuer         string
}

type RedisConfig struct {
	Addr string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", DefaultServerHost),
			Port:         getEnvInt("SERVER_PORT", DefaultServerPort),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", DefaultServerReadTimeout),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", DefaultServerWriteTimeout),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DATABASE_DSN", ""),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", DefaultLogLevel),
		},
		JWT: JWTConfig{
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "keys/jwt_rsa"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "keys/jwt_rsa.pub"),
			AccessTTL:      getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL:     getEnvDuration("JWT_REFRESH_TTL", 168*time.Hour),
			MFATempTTL:     getEnvDuration("JWT_MFA_TEMP_TTL", 5*time.Minute),
			Issuer:         getEnv("JWT_ISSUER", "observer"),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

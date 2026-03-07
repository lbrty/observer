package config

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultServerHost         = "localhost"
	DefaultServerPort         = 9000
	DefaultServerReadTimeout  = 30 * time.Second
	DefaultServerWriteTimeout = 30 * time.Second
	DefaultLogLevel           = "info"
	DefaultCookieMaxAge       = 2 * time.Hour
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Log       LogConfig
	JWT       JWTConfig
	Swagger   SwaggerConfig
	CORS      CORSConfig
	Cookie    CookieConfig
	RateLimit RateLimitConfig
	Storage   StorageConfig
}

type StorageConfig struct {
	Path string
}

type RedisConfig struct {
	URL string
}

type RateLimitConfig struct {
	LoginRate    int
	RegisterRate int
}

type CORSConfig struct {
	Origins []string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
	MaxAge   time.Duration
}

// HTTPSameSite converts the string config to http.SameSite.
func (c CookieConfig) HTTPSameSite() http.SameSite {
	switch strings.ToLower(c.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

type SwaggerConfig struct {
	Enabled bool
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
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
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
		Swagger: SwaggerConfig{
			Enabled: getEnvBool("SWAGGER_ENABLED", false),
		},
		CORS: CORSConfig{
			Origins: getEnvList("CORS_ORIGINS", []string{"http://localhost:5173"}),
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", ""),
			Secure:   getEnvBool("COOKIE_SECURE", true),
			SameSite: getEnv("COOKIE_SAME_SITE", "lax"),
			MaxAge:   getEnvDuration("COOKIE_MAX_AGE", DefaultCookieMaxAge),
		},
		RateLimit: RateLimitConfig{
			LoginRate:    getEnvInt("RATE_LIMIT_LOGIN", 10),
			RegisterRate: getEnvInt("RATE_LIMIT_REGISTER", 5),
		},
		Storage: StorageConfig{
			Path: getEnv("STORAGE_PATH", "data/uploads"),
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

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func getEnvList(key string, def []string) []string {
	if v := os.Getenv(key); v != "" {
		return strings.Split(v, ",")
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

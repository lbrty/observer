package config

import (
	"net/http"
	"strings"
	"time"

	"github.com/sultaniman/env"
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
	DevMode   bool
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
	Sentry    SentryConfig
}

type SentryConfig struct {
	DSN              string
	TracesSampleRate float64
}

func (c SentryConfig) Enabled() bool {
	return c.DSN != ""
}

type StorageConfig struct {
	Path        string // STORAGE_PATH,    default "data/uploads"
	Backend     string // STORAGE_BACKEND, default "local"
	S3Endpoint  string // S3_ENDPOINT,     default "" (uses AWS default resolver)
	S3Bucket    string // S3_BUCKET
	S3Region    string // S3_REGION,       default "us-east-1"
	S3AccessKey string // S3_ACCESS_KEY    (optional — falls back to SDK credential chain)
	S3SecretKey string // S3_SECRET_KEY    (optional — falls back to SDK credential chain)
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
	devMode := env.GetBool("DEV_MODE")
	return &Config{
		DevMode: devMode,
		Server: ServerConfig{
			Host:         strOr("SERVER_HOST", DefaultServerHost),
			Port:         intOr("SERVER_PORT", DefaultServerPort),
			ReadTimeout:  durOr("SERVER_READ_TIMEOUT", DefaultServerReadTimeout),
			WriteTimeout: durOr("SERVER_WRITE_TIMEOUT", DefaultServerWriteTimeout),
		},
		Database: DatabaseConfig{
			DSN: env.GetString("DATABASE_DSN"),
		},
		Redis: RedisConfig{
			URL: strOr("REDIS_URL", "redis://localhost:6379/0"),
		},
		Log: LogConfig{
			Level: strOr("LOG_LEVEL", DefaultLogLevel),
		},
		JWT: JWTConfig{
			PrivateKeyPath: strOr("JWT_PRIVATE_KEY_PATH", "keys/jwt_rsa"),
			PublicKeyPath:  strOr("JWT_PUBLIC_KEY_PATH", "keys/jwt_rsa.pub"),
			AccessTTL:      durOr("JWT_ACCESS_TTL", 30*time.Minute),
			RefreshTTL:     durOr("JWT_REFRESH_TTL", 168*time.Hour),
			MFATempTTL:     durOr("JWT_MFA_TEMP_TTL", 5*time.Minute),
			Issuer:         strOr("JWT_ISSUER", "observer"),
		},
		Swagger: SwaggerConfig{
			Enabled: env.GetBool("SWAGGER_ENABLED"),
		},
		CORS: CORSConfig{
			Origins: listOr("CORS_ORIGINS", []string{"http://localhost:5173"}),
		},
		Cookie: CookieConfig{
			Domain:   env.GetString("COOKIE_DOMAIN"),
			Secure:   boolOr("COOKIE_SECURE", !devMode),
			SameSite: strOr("COOKIE_SAME_SITE", "lax"),
			MaxAge:   durOr("COOKIE_MAX_AGE", DefaultCookieMaxAge),
		},
		RateLimit: RateLimitConfig{
			LoginRate:    intOr("RATE_LIMIT_LOGIN", 10),
			RegisterRate: intOr("RATE_LIMIT_REGISTER", 5),
		},
		Storage: StorageConfig{
			Path:        strOr("STORAGE_PATH", "data/uploads"),
			Backend:     strOr("STORAGE_BACKEND", "local"),
			S3Endpoint:  env.GetString("S3_ENDPOINT"),
			S3Bucket:    env.GetString("S3_BUCKET"),
			S3Region:    strOr("S3_REGION", "us-east-1"),
			S3AccessKey: env.GetString("S3_ACCESS_KEY"),
			S3SecretKey: env.GetString("S3_SECRET_KEY"),
		},
		Sentry: SentryConfig{
			DSN:              env.GetString("SENTRY_DSN"),
			TracesSampleRate: float64Or("SENTRY_TRACES_SAMPLE_RATE", 0.1),
		},
	}, nil
}

func strOr(key, def string) string {
	if v := env.GetString(key); v != "" {
		return v
	}
	return def
}

func intOr(key string, def int) int {
	if v := env.GetInt(key); v != 0 {
		return v
	}
	return def
}

// boolOr is needed only when the default is true (env.GetBool returns false on missing).
func boolOr(key string, def bool) bool {
	if v, err := env.GetBoolE(key); err == nil {
		return v
	}
	return def
}

func float64Or(key string, def float64) float64 {
	if v := env.GetFloat64(key); v != 0 {
		return v
	}
	return def
}

func durOr(key string, def time.Duration) time.Duration {
	if v := env.GetString(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func listOr(key string, def []string) []string {
	if v := env.GetString(key); v != "" {
		return strings.Split(v, ",")
	}
	return def
}

package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Addr        string
	DatabaseURL string
	JWTSecret   string

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	CookieDomain string
	CookieSecure bool
}

func MustLoad() Config {
	addr := getenv("APP_ADDR", ":8080")
	dbURL := mustEnv("DATABASE_URL")
	secret := mustEnv("JWT_SECRET")

	accessTTL := mustDuration("ACCESS_TTL", 15*time.Minute)
	refreshTTL := mustDuration("REFRESH_TTL", 7*24*time.Hour)

	return Config{
		Addr:         addr,
		DatabaseURL:  dbURL,
		JWTSecret:    secret,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
		CookieDomain: getenv("COOKIE_DOMAIN", "localhost"),
		CookieSecure: getenv("COOKIE_SECURE", "false") == "true",
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func mustDuration(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("invalid duration %s=%q: %v", k, v, err)
	}
	return d
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	GATEWAY_PORT string
	ORDER_PORT string
	// TIMELINE_PORT string

	ORDER_HOST string
	// TIMELINE_HOST string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.GATEWAY_PORT = cast.ToString(coalesce("GATEWAY_PORT", ":8010"))
	config.ORDER_PORT = cast.ToString(coalesce("ORDER_PORT", ":50060"))
	// config.TIMELINE_PORT = cast.ToString(coalesce("TIMELINE_PORT", ":50052"))
	config.ORDER_HOST = cast.ToString(coalesce("ORDER_HOST", "localhost"))
	// config.TIMELINE_HOST = cast.ToString(coalesce("TIMELINE_HOST", "timeline-service"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}

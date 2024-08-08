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
	COURIER_PORT string

	ORDER_HOST string
	COURIER_HOST string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.GATEWAY_PORT = cast.ToString(coalesce("GATEWAY_PORT", ":8010"))
	config.ORDER_PORT = cast.ToString(coalesce("ORDER_PORT", ":50060"))
	config.COURIER_PORT = cast.ToString(coalesce("COURIER_PORT", ":50061"))
	config.ORDER_HOST = cast.ToString(coalesce("ORDER_HOST", "localhost"))
	config.COURIER_HOST = cast.ToString(coalesce("COURIER_HOST", "localhost"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}

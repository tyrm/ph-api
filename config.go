package main

import (
	"fmt"
	"os"
)

type Config struct {
	DBEngine        string

	HTTPCorsOrigin  string

	RedisAddr       string
	RedisPassword   string
}

func CollectConfig() (config Config) {
	var missingEnv []string

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// HTTP_CORS_ORIGIN
	config.HTTPCorsOrigin = os.Getenv("HTTP_CORS_ORIGIN")

	// REDIS_ADDR
	var envRedisAddress string = os.Getenv("REDIS_ADDR")

	if envRedisAddress == "" {
		config.RedisAddr = "localhost:6379"
	} else {
		config.RedisAddr = envRedisAddress
	}

	// REDIS_PASSWORD
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")

	// Validation
	if len(missingEnv) > 0 {
		var msg string = fmt.Sprintf("Environment variables missing: %v", missingEnv)
		logger.Criticalf(msg)
		panic(fmt.Sprint(msg))
	}

	return
}

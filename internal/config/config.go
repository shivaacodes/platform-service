package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	Env           string `mapstructure:"ENV"`
}

var C Config

// type Config defines structured fields for app settings (PORT, REDIS_ADDR, etc.).

// viper.SetConfigFile(".env") tells Viper to read from a local .env file if present.

// viper.AutomaticEnv() makes environment variables override values from .env.

// viper.SetDefault(...) provides fallback values if nothing else is set.

// viper.ReadInConfig() tries to read .env. If missing, it just logs a message.

// viper.Unmarshal(&C) maps all loaded settings into the global Config struct C.

func Load() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env") // treat as env format
	viper.AddConfigPath(".")   // search in current directory
	viper.AutomaticEnv()       // allow OS environment override

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("ENV", "development")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("no .env file found, using defaults and environment vars")
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("config load error: %v", err)
	}
}

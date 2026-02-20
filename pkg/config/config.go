package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Notification NotificationConfig `mapstructure:"notification"`
	RateLimit    RateLimitConfig    `mapstructure:"rate_limit"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

type NotificationConfig struct {
	FirebaseCredentialsFile string `mapstructure:"firebase_credentials_file"`
}

type RateLimitConfig struct {
	Global int `mapstructure:"global"`
	Auth   int `mapstructure:"auth"`
	Public int `mapstructure:"public"`
	Sync   int `mapstructure:"sync"`
}

func LoadConfig() *Config {
	viper.SetConfigFile("config.yaml")
	viper.AutomaticEnv()

	// Explicit bindings for Docker Env Vars
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("server.port", "PORT")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using defaults based on env or structs")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Printf("Loaded Config: %+v", config) // Debug Log
	return &config
}

package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from environment / .env file.
type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	FCM      FCMConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret         string
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
}

type StorageConfig struct {
	Type      string
	Path      string
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type FCMConfig struct {
	ServerKey string
}

type CORSConfig struct {
	Origins string
}

// Load reads configuration from environment variables / .env file.
func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config: no .env file found, reading from environment: %v", err)
	}

	accessExpiry, _ := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRY"))
	if accessExpiry == 0 {
		accessExpiry = 15 * time.Minute
	}
	refreshExpiry, _ := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour
	}

	return &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
			Name: viper.GetString("APP_NAME"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			Timezone: viper.GetString("DB_TIMEZONE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Storage: StorageConfig{
			Type:      viper.GetString("STORAGE_TYPE"),
			Path:      viper.GetString("STORAGE_PATH"),
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		},
		FCM: FCMConfig{
			ServerKey: viper.GetString("FCM_SERVER_KEY"),
		},
		CORS: CORSConfig{
			Origins: viper.GetString("CORS_ORIGINS"),
		},
	}
}

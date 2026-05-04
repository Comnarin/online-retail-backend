package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LIFF     LIFFConfig
	MinIO    MinIOConfig
}

type AppConfig struct {
	Port string `mapstructure:"PORT"`
	Env  string `mapstructure:"APP_ENV"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"DATABASE_URL"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"JWT_SECRET"`
	ExpiryHours     int    `mapstructure:"JWT_EXPIRY_HOURS"`
	RefreshDays     int    `mapstructure:"JWT_REFRESH_DAYS"`
}

type LIFFConfig struct {
	ID        string `mapstructure:"LIFF_ID"`
	ChannelID string `mapstructure:"LINE_CHANNEL_ID"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket    string `mapstructure:"MINIO_BUCKET"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_DAYS", 30)
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("LIFF_ID", "PLACEHOLDER_LIFF_ID")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET", "retail")
	viper.SetDefault("MINIO_USE_SSL", false)

	_ = viper.ReadInConfig() // won't fail if .env not found

	cfg := &Config{}
	cfg.App.Port = viper.GetString("PORT")
	cfg.App.Env = viper.GetString("APP_ENV")
	cfg.Database.DSN = viper.GetString("DATABASE_URL")
	cfg.Redis.Addr = viper.GetString("REDIS_ADDR")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = viper.GetInt("REDIS_DB")
	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.ExpiryHours = viper.GetInt("JWT_EXPIRY_HOURS")
	cfg.JWT.RefreshDays = viper.GetInt("JWT_REFRESH_DAYS")
	cfg.LIFF.ID = viper.GetString("LIFF_ID")
	cfg.LIFF.ChannelID = viper.GetString("LINE_CHANNEL_ID")
	cfg.MinIO.Endpoint = viper.GetString("MINIO_ENDPOINT")
	cfg.MinIO.AccessKey = viper.GetString("MINIO_ACCESS_KEY")
	cfg.MinIO.SecretKey = viper.GetString("MINIO_SECRET_KEY")
	cfg.MinIO.Bucket = viper.GetString("MINIO_BUCKET")
	cfg.MinIO.UseSSL = viper.GetBool("MINIO_USE_SSL")

	return cfg, nil
}


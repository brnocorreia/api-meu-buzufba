package config

import (
	"log"
	"sync"

	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	AppName     string `mapstructure:"APP_NAME"`
	DebugMode   bool   `mapstructure:"DEBUG"`

	PostgresDSN string `mapstructure:"DB_POSTGRES_DSN"`
	ResendKey   string `mapstructure:"RESEND_API_KEY"`

	JWTSecretKey            string `mapstructure:"JWT_SECRET"`
	JWTAccessTokenDuration  string `mapstructure:"JWT_ACCESS_TOKEN_DURATION"`
	JWTRefreshTokenDuration string `mapstructure:"JWT_REFRESH_TOKEN_DURATION"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       string `mapstructure:"REDIS_DB"`
}

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			logging.Error("error reading config file", err,
				zap.String("journey", "config"))
			log.Fatalf("error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			logging.Error("error unmarshalling config", err,
				zap.String("journey", "config"))
			log.Fatalf("error unmarshalling config, %s", err)
		}
	})

	return config
}

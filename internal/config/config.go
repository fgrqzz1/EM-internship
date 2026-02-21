package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Name       string `mapstructure:"name"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	MaxOpenCon int    `mapstructure:"max_open_con"`
	MaxIdleCon int    `mapstructure:"max_idle_con"`
}

type LoggingConfig struct {
	Level       string `mapstructure:"level"`
	Development bool   `mapstructure:"development"`
}

func LoadConfig(logger *zap.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	//viper.AddConfigPath(".") // при расположения yaml в корне

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("Config file not found", zap.Error(err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal , %v", err)
	}

	return &cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	return os.ExpandEnv(dsn)
}

func NewLogger(cfg *LoggingConfig) (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()

	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapCfg.Level.SetLevel(level)

	return zapCfg.Build()
}

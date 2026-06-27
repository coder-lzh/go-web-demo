package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name    string `mapstructure:"name" json:"name"`
	Version string `mapstructure:"version" json:"version"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Mode string `mapstructure:"mode" json:"mode"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host" json:"host"`
	Port            int    `mapstructure:"port" json:"port"`
	Username        string `mapstructure:"username" json:"username"`
	Password        string `mapstructure:"password" json:"-"`
	Database        string `mapstructure:"database" json:"database"`
	Charset         string `mapstructure:"charset" json:"charset"`
	Loc             string `mapstructure:"loc" json:"loc"`
	LogMode         bool   `mapstructure:"log_mode" json:"log_mode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Format     string `mapstructure:"format" json:"format"` // json, text
	FilePath   string `mapstructure:"file_path" json:"file_path"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
	Console    bool   `mapstructure:"console" json:"console"`
}

type Config struct {
	Env      string         `json:"env"`
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
}

func LoadConfig(env string) (*Config, error) {
	if env == "" {
		env = os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
	}

	v := viper.New()
	v.SetConfigType("yaml")

	v.AddConfigPath("./configs")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	envConfigFile := fmt.Sprintf("config.%s", env)
	v.SetConfigName(envConfigFile)
	if err := v.MergeInConfig(); err != nil {
		fmt.Printf("Warning: configs/%s.yaml not found, skipping environment config\n", envConfigFile)
	}

	v.SetConfigName("config.local")
	if err := v.MergeInConfig(); err == nil {
		fmt.Println("Loaded local config override")
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.Env = env

	return &config, nil
}

func (c *Config) IsDev() bool {
	return c.Env == "dev" || c.Env == "development"
}

func (c *Config) IsProd() bool {
	return c.Env == "prod" || c.Env == "production"
}

func (c *Config) IsTest() bool {
	return c.Env == "test" || c.Env == "testing"
}

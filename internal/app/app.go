package app

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"go-web-demo/internal/config"
	"go-web-demo/internal/db"
)

func InitConfig(env string) (*config.Config, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	setLogLevel(cfg.Log.Level)

	hlog.Infof("Starting %s v%s in %s mode...", cfg.App.Name, cfg.App.Version, cfg.Env)

	return cfg, nil
}

func InitDatabase(cfg *config.Config) error {
	if err := db.InitDB(cfg); err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	hlog.Info("Database connection established successfully")
	return nil
}

func CreateServer(cfg *config.Config) *server.Hertz {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return server.Default(server.WithHostPorts(addr))
}

func setLogLevel(level string) {
	switch level {
	case "trace":
		hlog.SetLevel(hlog.LevelTrace)
	case "debug":
		hlog.SetLevel(hlog.LevelDebug)
	case "info":
		hlog.SetLevel(hlog.LevelInfo)
	case "warn":
		hlog.SetLevel(hlog.LevelWarn)
	case "error":
		hlog.SetLevel(hlog.LevelError)
	default:
		hlog.SetLevel(hlog.LevelInfo)
	}
}
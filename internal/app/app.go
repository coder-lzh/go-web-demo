package app

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"

	"go-web-demo/internal/config"
	"go-web-demo/internal/db"
	"go-web-demo/internal/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitConfig(env string) (*config.Config, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := initLogger(cfg); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	logger.Infof("Starting %s v%s in %s mode...", cfg.App.Name, cfg.App.Version, cfg.Env)

	return cfg, nil
}

func InitDatabase(cfg *config.Config) error {
	if err := db.InitDB(cfg); err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	logger.Info("Database connection established successfully")
	return nil
}

func CreateServer(cfg *config.Config) *server.Hertz {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return server.Default(server.WithHostPorts(addr))
}

func initLogger(cfg *config.Config) error {
	logCfg := &logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		Console:    cfg.Log.Console,
	}
	return logger.Init(logCfg)
}

func InitTracer() {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	logger.Info("OpenTelemetry tracer initialized")
}

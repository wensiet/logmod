package logmod

import (
	"fmt"
	slogmulti "github.com/samber/slog-multi"
	"github.com/wensiet/logmod/loki"
	"github.com/wensiet/logmod/slogloki"
	"log/slog"
	"os"
)

type Options struct {
	Env     string
	Service string
	Loki    struct {
		Host string
		Port int
	}
}

func New(opts Options) *slog.Logger {
	lokiClientError := false
	lokiConfig, err := loki.NewDefaultConfig(
		fmt.Sprintf("http://%s:%d/loki/api/v1/push", opts.Loki.Host, opts.Loki.Port),
	)

	lokiConfig.TenantID = "xyz"
	if err != nil {
		lokiClientError = true
	}

	lokiClient, err := loki.New(lokiConfig)
	if err != nil {
		lokiClientError = true
	}

	var logger *slog.Logger

	switch opts.Env {
	case "local":
		logger = slog.New(
			slogmulti.Fanout(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
				func() slog.Handler {
					if !lokiClientError {
						return slogloki.Option{Level: slog.LevelDebug, Client: lokiClient}.NewLokiHandler()
					} else {
						return nil
					}
				}(),
			),
		).With("environment", "local")
	case "production":
		logger = slog.New(
			slogmulti.Fanout(
				func() slog.Handler {
					if !lokiClientError {
						return slogloki.Option{Level: slog.LevelInfo, Client: lokiClient}.NewLokiHandler()
					} else {
						return nil
					}
				}(),
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			),
		).With("environment", "production")
	case "test":
		logger = slog.New(
			slogmulti.Fanout(
				func() slog.Handler {
					if !lokiClientError {
						return slogloki.Option{Level: slog.LevelDebug, Client: lokiClient}.NewLokiHandler()
					} else {
						return nil
					}
				}(),
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		).With("environment", "tests")
	}

	logger = logger.With("service", opts.Service)

	return logger
}

package myfoodserver

import (
	"context"
	"fmt"
	"net/http"

	handler "github.com/devldavydov/myfood/internal/myfoodserver/handlers"
	"github.com/devldavydov/myfood/internal/myfoodserver/templates"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service struct {
	settings *ServerSettings
	logger   *zap.Logger
	stg      storage.Storage
}

func NewService(settings *ServerSettings, logger *zap.Logger) (*Service, error) {
	return &Service{settings: settings, logger: logger}, nil
}

func (r *Service) Run(ctx context.Context) error {
	// Init HTTP API
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.HTMLRender = templates.NewTemplateRenderer()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	handler.Init(router, r.stg, r.logger)

	// Start server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", r.settings.RunAddress.Hostname(), r.settings.RunAddress.Port()),
		Handler: router,
	}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		return fmt.Errorf("service exited with err: %w", err)
	case <-ctx.Done():
		r.logger.Info("Service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), r.settings.ShutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("service shutdown err: %w", err)
		}

		r.logger.Info("Service finished")
		return nil
	}
}

package myfoodserver

import (
	"context"
	"fmt"
	"net/http"

	handler "github.com/devldavydov/myfood/internal/myfoodserver/handlers"
	"github.com/devldavydov/myfood/internal/myfoodserver/templates"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const SessionName = "MyFoodSession"

type Service struct {
	settings *ServerSettings
	logger   *zap.Logger
	stg      storage.Storage
}

func NewService(settings *ServerSettings, logger *zap.Logger) (*Service, error) {
	stg, err := storage.NewStorageSQLite(settings.DBFilePath)
	if err != nil {
		return nil, err
	}

	return &Service{settings: settings, stg: stg, logger: logger}, nil
}

func (r *Service) Run(ctx context.Context) error {
	// Init HTTP API
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.HTMLRender = templates.NewTemplateRenderer()

	store := cookie.NewStore([]byte(r.settings.SessionSecret))

	router.Use(
		gzip.Gzip(gzip.DefaultCompression),
		sessions.Sessions(SessionName, store),
	)

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

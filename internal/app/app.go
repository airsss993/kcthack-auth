package app

import (
	"log"

	"github.com/kcthack-auth/internal/config"
	"github.com/kcthack-auth/internal/handler"
	"github.com/kcthack-auth/internal/repository"
	"github.com/kcthack-auth/internal/serser"
	"github.com/kcthack-auth/internal/service"
	"github.com/kcthack-auth/pkg/auth"
	"github.com/kcthack-auth/pkg/database"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %e", err)
	}

	db := database.ConnDB(cfg)
	authRepo := repository.NewAuthRepo(db)
	sessRepo := repository.NewSessionRepo(db)
	tm := auth.NewManager(cfg.JWT.JWTSecret)
	authService := service.NewAuthService(authRepo, sessRepo, tm, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL)
	services := service.NewServices(*authService)
	newHandler := handler.NewHandler(services)
	r := newHandler.Init(cfg)

	server := serser.NewServer(r, *cfg)
	if err := server.Start(); err != nil {
		log.Fatal("failed to start http server")
	}
}

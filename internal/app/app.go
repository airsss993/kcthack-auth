package app

import (
	"log"

	"github.com/kcthack-auth/internal/config"
	"github.com/kcthack-auth/internal/handler"
	"github.com/kcthack-auth/pkg/database"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %e", err)
	}

	_ = database.ConnDB(cfg)

	r := handler.NewRouter()
	r.Run(cfg.HTTP.Port)
}

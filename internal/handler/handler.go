package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kcthack-auth/internal/config"
	v1 "github.com/kcthack-auth/internal/handler/v1"
	"github.com/kcthack-auth/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	h.initAPI(r)

	return r
}

func (h *Handler) initAPI(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		v1Group := api.Group("/v1")
		{
			handlerV1 := v1.NewHandler(*h.services)
			handlerV1.Init(v1Group)
		}

	}

}

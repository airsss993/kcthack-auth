package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kcthack-auth/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services service.Services) *Handler {
	return &Handler{services: &services}
}

func (h *Handler) Init(a *gin.RouterGroup) {
	user := a.Group("/user")
	{
		user.POST("/register", h.register)
		user.POST("/login", h.login)
		user.POST("/logout", h.logout)
		user.POST("/refresh", h.refresh)
	}
}

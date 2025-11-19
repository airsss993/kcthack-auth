package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kcthack-auth/internal/service"
)

type userRegisterReq struct {
	FirstName string `json:"first_name" binding:"required,max=32"`
	LastName  string `json:"last_name" binding:"required,max=32"`
	Email     string `json:"email" binding:"required,email,max=32"`
	Password  string `json:"password" binding:"required,min=6"`
}

type userLoginReq struct {
	Email    string `json:"email" binding:"required,email,max=32"`
	Password string `json:"password" binding:"required,min=3"`
}

type userLoginResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (h *Handler) register(c *gin.Context) {
	var req userRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.services.AuthService.Register(c.Request.Context(), service.RegisterReq{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", resp.AccessToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "registration successful",
	})
}

func (h *Handler) login(c *gin.Context) {
	var req userLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.services.AuthService.Login(c.Request.Context(), &service.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", resp.AccessToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}

func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token not found in cookies",
		})
		return
	}

	if err := h.services.AuthService.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token not found in cookies",
		})
		return
	}

	resp, err := h.services.AuthService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", resp.AccessToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(resp.ExpiresAt.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "tokens refreshed successfully",
	})
}

package handlers

import (
	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Elimina esta función y la registración del middleware
//func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de autenticación no proporcionado"})
//			return
//		}
//
//		// Extraer el token del encabezado
//		parts := strings.Split(authHeader, " ")
//		if len(parts) != 2 || parts[0] != "Bearer" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
//			return
//		}
//
//		token := parts[1]
//		userID, err := h.authService.ValidateToken(token)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
//			return
//		}
//
//		c.Set("userID", userID)
//		c.Next()
//	}
//}

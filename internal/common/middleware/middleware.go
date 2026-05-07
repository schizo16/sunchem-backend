package middleware

import (
	"net/http"
	"strings"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/utils"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		claims, err := utils.ParseToken(secret, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminOnly(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, errors.ErrForbidden)
			return
		}
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.HttpStatus, appErr)
			} else {
				c.JSON(500, errors.ErrInternalServer)
			}
		}
	}
}

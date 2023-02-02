package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phone_numbers_checker/utils/token"
)

func AuthRequired(c *gin.Context) {
	err := token.IsTokenValid(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}

package middlewares

import (
	"net/http"

	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/auth"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	errList := make(map[string]string)
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

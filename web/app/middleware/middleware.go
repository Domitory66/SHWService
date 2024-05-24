package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Param("Auth")
		log.Println(c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

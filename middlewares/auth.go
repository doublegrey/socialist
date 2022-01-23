package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		if id == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Code": http.StatusUnauthorized, "Error": "unauthorized", "Event": uuid.NewString()})
			c.Abort()
		}
	}
}

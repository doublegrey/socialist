package api

import (
	"log"

	"github.com/doublegrey/socialist/api/handlers"
	"github.com/doublegrey/socialist/datastore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	sessionsStore, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("Failed to create sessions store: %v\n", err)
	}
	sessionsStore.Options(sessions.Options{MaxAge: 86400})
	r.Use(sessions.Sessions("session", sessionsStore))

	authHandler := handlers.NewAuthHandler(datastore.Connections)

	apiGroup := r.Group("api")
	authGroup := apiGroup.Group("auth")

	authGroup.POST("login", authHandler.Login)
	authGroup.POST("register", authHandler.Register)
	authGroup.GET("logout", authHandler.Logout)
}

package main

import (
	"github.com/doublegrey/socialist/api"
	"github.com/doublegrey/socialist/datastore"
	"github.com/doublegrey/socialist/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	datastore.New()
	r := gin.Default()
	r.Use(middlewares.CorsAllowAll())
	api.NewRouter(r)
	r.Run(":8000")
}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/routes"
)

func main() {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	routes.SessionRoutes(engine)
	// Listen and Server in 0.0.0.0:8080
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}

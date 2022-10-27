package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xgss/routes"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	routes.SessionRoutes(engine)
	err := engine.Run("0.0.0.0:42100")
	if err != nil {
		return
	}
}

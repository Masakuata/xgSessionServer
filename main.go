package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"xgss/routes"
)

func main() {
	var engine *gin.Engine

	if _, isDebug := os.LookupEnv("DEBUG"); isDebug {
		engine = gin.Default()
		gin.SetMode(gin.DebugMode)
	} else {
		engine = gin.New()
		gin.SetMode(gin.ReleaseMode)
	}

	engine.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	routes.SessionRoutes(engine)
	routes.RequireRoutes(engine)
	err := engine.Run("0.0.0.0:42100")
	if err != nil {
		return
	}
}

package server

import (
	"github.com/Dpbm/quantumRestAPI/db"
	"github.com/Dpbm/quantumRestAPI/middlewares"
	routes "github.com/Dpbm/quantumRestAPI/routes"
	"github.com/gin-gonic/gin"
)

func SetupServer(dbInstance *db.DB) *gin.Engine {
	server := gin.Default()

	server.Use(middlewares.DB(dbInstance))

	server.GET("/job/:id", routes.GetJob)
	server.POST("/plugin/:name", routes.AddPlugin)

	return server
}

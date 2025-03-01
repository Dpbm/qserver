package server

import (
	"github.com/Dpbm/quantumRestAPI/db"
	docs "github.com/Dpbm/quantumRestAPI/docs"
	"github.com/Dpbm/quantumRestAPI/middlewares"
	routes "github.com/Dpbm/quantumRestAPI/routes"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupServer(dbInstance *db.DB) *gin.Engine {
	server := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	server.Use(middlewares.DB(dbInstance))

	v1 := server.Group("/api/v1")
	{
		job := v1.Group("/job")
		{
			job.GET("/result/:id", routes.GetJobResult)
			job.DELETE("/:id", routes.DeleteJob)
			job.GET("/:id", routes.GetJob)
		}

		jobs := v1.Group("/jobs")
		{
			jobs.GET("/", routes.GetJobs)
		}

		plugin := v1.Group("/plugin")
		{
			plugin.POST("/:name", routes.AddPlugin)
		}
	}

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return server
}

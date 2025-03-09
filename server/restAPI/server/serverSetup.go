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

	// some routes, for a reason that I don't know
	// redirect the request to the same page if a
	// path parameter isn't passed. To solve that,
	// I forced a / path to return a 404 error

	v1 := server.Group("/api/v1")
	{
		job := v1.Group("/job")
		{
			job.GET("/", routes.Page404)
			job.GET("/result/", routes.Page404)

			job.GET("/:id", routes.GetJob)
			job.GET("/result/:id", routes.GetJobResult)
			job.PUT("/cancel/:id", routes.CancelJob)
			job.DELETE("/:id", routes.DeleteJob)
		}

		jobs := v1.Group("/jobs")
		{
			jobs.GET("/", routes.GetJobs)
		}

		history := v1.Group("/history")
		{
			history.GET("/", routes.GetHistory)
		}

		plugin := v1.Group("/plugin")
		{
			plugin.POST("/:name", routes.AddPlugin)
			plugin.DELETE("/:name", routes.DeletePlugin)
		}

		backend := v1.Group("/backend")
		{
			backend.GET("/", routes.Page404)

			backend.GET("/:name", routes.GetBackend)
		}

		backends := v1.Group("/backends")
		{
			backends.GET("/", routes.GetBackends)
		}

	}

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return server
}

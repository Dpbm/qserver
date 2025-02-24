package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Dpbm/quantumRestAPI/db"
	logger "github.com/Dpbm/quantumRestAPI/log"
	"github.com/Dpbm/quantumRestAPI/middlewares"
	"github.com/Dpbm/quantumRestAPI/routes"
	"github.com/Dpbm/quantumRestAPI/types"
)

func main() {
	port := os.Getenv("PORT")
	if !types.ValidIntFromEnv(port) {
		logger.LogFatal(errors.New("invalid Server Port"))
		os.Exit(1) // ensure the program is going to exit
	}

	dbInstance := db.DB{}
	dbInstance.Connect(&db.Postgres{}) // on error it should exit the program with return code 1

	server := gin.Default()

	server.Use(middlewares.DB(&dbInstance))

	server.GET("/job/:id", routes.GetJob)
	server.POST("/plugin/:name", routes.AddPlugin)

	portString := fmt.Sprintf(":%s", port)
	server.Run(portString)
}

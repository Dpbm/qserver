package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Dpbm/quantumRestAPI/db"
	logger "github.com/Dpbm/quantumRestAPI/log"
	"github.com/Dpbm/quantumRestAPI/server"
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
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	portString := fmt.Sprintf(":%s", port)
	server.Run(portString)
}

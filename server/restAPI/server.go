package main

import (
	"fmt"
	"os"

	"github.com/Dpbm/quantumRestAPI/db"
	dbDefinition "github.com/Dpbm/shared/db"

	"github.com/Dpbm/quantumRestAPI/server"
	"github.com/Dpbm/shared/format"
)

func main() {
	port := format.PortEnvToInt(os.Getenv("PORT")) // it must execute os.Exit(1) if the port is invalid

	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Postgres{}) // on error it should exit the program with return code 1
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	portString := fmt.Sprintf(":%d", port)
	server.Run(portString)
}

package format

import (
	"errors"
	"os"
	"strconv"

	logger "github.com/Dpbm/quantumRestAPI/log"
)

func PortEnvToInt(env string) int {
	port, err := strconv.Atoi(env)

	if err != nil || port < 0 {
		logger.LogFatal(errors.New("failed on convert port env to int"))
		os.Exit(1) // just to ensure the program is going to close
	}

	return port
}

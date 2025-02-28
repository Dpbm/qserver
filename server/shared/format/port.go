package format

import (
	"errors"
	"os"
	"strconv"

	logger "github.com/Dpbm/shared/log"
)

func PortEnvToInt(env string) uint32 {
	port, err := strconv.Atoi(env)

	if err != nil || port < 0 {
		logger.LogFatal(errors.New("failed on convert port env to int"))
		os.Exit(1) // just to ensure the program is going to close
	}

	return uint32(port)
}

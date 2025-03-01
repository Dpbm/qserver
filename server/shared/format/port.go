package format

import (
	"errors"
	"os"
	"strconv"

	logger "github.com/Dpbm/shared/log"
)

func PortEnvToInt(env string) uint32 {
	port, err := StrToUint(env)

	if err != nil {
		logger.LogFatal(errors.New("failed on convert port env to int"))
		os.Exit(1) // just to ensure the program is going to close
	}

	return port
}

func StrToUint(value string) (uint32, error) {
	uintValue, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	if uintValue < 0 {
		return 0, errors.New("invalid uint string")
	}

	return uint32(uintValue), nil

}

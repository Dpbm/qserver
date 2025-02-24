package types

import (
	"github.com/Dpbm/quantumRestAPI/format"
)

func ValidIntFromEnv(env string) bool {
	format.PortEnvToInt(env)

	// once format.PortEnvToInt exits if the value is invalid,
	// we can simply return true once the program hasn't exited
	return true

}

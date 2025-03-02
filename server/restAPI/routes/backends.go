package routes

import (
	"errors"

	"github.com/Dpbm/quantumRestAPI/types"
	"github.com/Dpbm/quantumRestAPI/utils"
	logger "github.com/Dpbm/shared/log"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @version 1.0
// @Summary get backend data
// @Schemes http
// @Description backend data by backend name
// @Tags backends
// @Param name path string true "Backend name"
// @Produce json
// @Success 200 {object} types.BackendData
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "It wasn't possible to find the backend"
// @Router /backend/{name} [get]
func GetBackend(context *gin.Context) {
	var backend types.BackendByName
	err := context.ShouldBindUri(&backend)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	result, err := db.GetBackend(backend.Name)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Could not find the required backend"})
		return
	}

	context.JSON(200, result)
}

package routes

import (
	"github.com/Dpbm/quantumRestAPI/types"
	"github.com/Dpbm/quantumRestAPI/utils"
	logger "github.com/Dpbm/shared/log"
	"github.com/gin-gonic/gin"
)

func AddPlugin(context *gin.Context) {
	var plugin types.AddPluginByName
	err := context.ShouldBindUri(&plugin)
	// TODO: test this part
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": err.Error()})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	pluginName := plugin.Name
	backends, err := utils.GetBackendsList(pluginName)

	if err != nil || len(*backends) <= 0 {
		logger.LogError(err)
		context.JSON(500, map[string]string{"msg": "Failed on get backends!"})
		return
	}

	err = db.SaveBackends(backends, pluginName)
	// TODO: test this part
	if err != nil {
		logger.LogError(err)
		context.JSON(500, map[string]string{"msg": "Failed on save data on DB!"})
		return
	}

	context.JSON(201, map[string]string{"msg": "added plugin"})
}

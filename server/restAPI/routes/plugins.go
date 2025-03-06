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
// @Summary add plugin
// @Schemes http
// @Description add plugin by name
// @Tags plugins
// @Param name path string true "Plugin name as shown in the github org"
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid Name parameter"
// @Failure 500 {object} map[string]string "Couldn't connect to database or get the plugin info from github"
// @Failure 404 {object} map[string]string "No results for this name"
// @Router /plugin/{name} [post]
func AddPlugin(context *gin.Context) {
	var plugin types.PluginByName
	err := context.ShouldBindUri(&plugin)
	// TODO: test this sort of error
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
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

// @BasePath /api/v1
// @version 1.0
// @Summary delete plugin
// @Schemes http
// @Description delete all data related to this plugin name
// @Tags plugins
// @Param name path string true "Plugin Name"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid Name parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "No results for this Name"
// @Router /plugin/{name} [delete]
func DeletePlugin(context *gin.Context) {
	var plugin types.PluginByName
	err := context.ShouldBindUri(&plugin)
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

	err = db.DeletePlugin(plugin.Name)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Failed on delete your plugin data. Remeber to wait your pending and running jobs to finished before deleting a plugin!"})
		return
	}

	context.JSON(200, map[string]string{"msg": "Success"})
}

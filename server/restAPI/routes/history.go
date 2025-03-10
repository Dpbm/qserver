package routes

import (
	"errors"

	"github.com/Dpbm/quantumRestAPI/utils"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @version 1.0
// @Summary get history data
// @Schemes http
// @Description get jobs history
// @Tags history
// @Param cursor query int false "Last id(pointer) gotten from db"
// @Produce json
// @Success 200 {object} []types.Historydata
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Router /history [get]
func GetHistory(context *gin.Context) {
	cursor := context.Query("cursor")
	cursorValue, err := format.StrToUint(cursor)

	if err != nil {
		logger.LogError(err)
		cursorValue = 0
	}

	db, ok := utils.GetDBFromContext(context)
	// TODO: TEST THIS PART
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	result, err := db.GetHistoryData(cursorValue)
	if err != nil {
		logger.LogError(err)
		context.JSON(200, map[string]any{})
		return
	}

	context.JSON(200, result)
}

package routes

import (
	"errors"

	"github.com/Dpbm/quantumRestAPI/types"
	"github.com/Dpbm/quantumRestAPI/utils"
	logger "github.com/Dpbm/shared/log"
	"github.com/gin-gonic/gin"
)

func GetJob(context *gin.Context) {
	var job types.GetJobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": err.Error()})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	result, err := db.GetJobData(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Results Data not found!"})
		return
	}

	context.JSON(200, result)
}

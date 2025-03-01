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
// @Summary get job results
// @Schemes http
// @Description get job results by ID
// @Tags jobs
// @Param id path string true "Job ID"
// @Produce json
// @Success 200 {object} types.JobData
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "No results for this ID"
// @Router /job/{id} [get]
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

// @BasePath /api/v1
// @version 1.0
// @Summary delete job data
// @Schemes http
// @Description delete all data related to this job id
// @Tags jobs
// @Param id path string true "Job ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "No results for this ID"
// @Router /job/{id} [delete]
func DeleteJob(context *gin.Context) {
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

	err = db.DeleteJobData(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Failed on delete your job data!"})
		return
	}

	context.JSON(200, map[string]string{"msg": "Sucess"})
}

// TODO: ADD CANCEL JOB (PUT)

// TODO: List all jobs
// TODO: list jobs by id

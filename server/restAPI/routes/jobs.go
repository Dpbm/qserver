package routes

import (
	"errors"

	"github.com/Dpbm/quantumRestAPI/types"
	"github.com/Dpbm/quantumRestAPI/utils"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @version 1.0
// @Summary get job data
// @Schemes http
// @Description get all data from job by ID
// @Tags jobs
// @Param id path string true "job ID"
// @Produce json
// @Success 200 {object} types.JobData
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "It wasn't possible to find the job"
// @Router /job/{id} [get]
func GetJob(context *gin.Context) {
	var job types.JobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	// TODO: TEST THIS PART
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	result, err := db.GetJob(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Could not find the required job"})
		return
	}

	context.JSON(200, result)
}

// @BasePath /api/v1
// @version 1.0
// @Summary get jobs data
// @Schemes http
// @Description get all data from jobs
// @Tags jobs
// @Param cursor query int false "Last id(pointer) gotten from db"
// @Produce json
// @Success 200 {object} []types.JobData
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Router /jobs [get]
func GetJobs(context *gin.Context) {
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

	result, err := db.GetJobsData(cursorValue)
	if err != nil {
		logger.LogError(err)
		context.JSON(200, map[any]any{})
		return
	}

	context.JSON(200, result)
}

// @BasePath /api/v1
// @version 1.0
// @Summary get job results
// @Schemes http
// @Description get job results by ID
// @Tags jobs
// @Param id path string true "Job ID"
// @Produce json
// @Success 200 {object} types.JobResultData
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection"
// @Failure 404 {object} map[string]string "No results for this ID or an issue occoured"
// @Router /job/result/{id} [get]
func GetJobResult(context *gin.Context) {
	var job types.JobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	// TODO: TEST THIS PART TOO
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	result, err := db.GetJobResult(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(404, map[string]string{"msg": "Results Data not found, Or a problem occured during the execution!"})
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
// @Failure 500 {object} map[string]string "Failed during DB connection or error during deletion"
// @Router /job/{id} [delete]
func DeleteJob(context *gin.Context) {
	var job types.JobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	// TODO: TEST THIS PART
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	err = db.DeleteJobData(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(500, map[string]string{"msg": "Failed on delete your job data!"})
		return
	}

	context.JSON(200, map[string]string{"msg": "Success"})
}

// @BasePath /api/v1
// @version 1.0
// @Summary cancel job
// @Schemes http
// @Description cancel a job before running it
// @Tags jobs
// @Param id path string true "Job ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 500 {object} map[string]string "Failed during DB connection or during db update"
// @Router /job/cancel/{id} [put]
func CancelJob(context *gin.Context) {
	var job types.JobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		logger.LogError(err)
		context.JSON(400, map[string]string{"msg": "Invalid Parameter"})
		return
	}

	db, ok := utils.GetDBFromContext(context)
	// TODO: TEST THIS PART
	if !ok || db == nil {
		logger.LogError(errors.New("failed on get DB from context"))
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
		return
	}

	err = db.CancelJob(job.ID)
	if err != nil {
		logger.LogError(err)
		context.JSON(500, map[string]string{"msg": "Failed on cancel your job!"})
		return
	}

	context.JSON(200, map[string]string{"msg": "Success"})
}

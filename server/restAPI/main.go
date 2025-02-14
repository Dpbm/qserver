package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type GetJobById struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func getJobData(jobId string, db *sql.DB) (sql.Result, error) {
	return db.Exec(`
		SELECT * 
		FROM results
		WHERE job_id=$1 
	`, (jobId))
}

func getJob(context *gin.Context) {
	var job GetJobById
	error := context.ShouldBindUri(&job)
	if error != nil {
		context.JSON(400, map[string]string{"msg": error.Error()})
		return
	}

	db, ok := context.MustGet("db").(*sql.DB)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
	}

	result, error := getJobData(job.ID, db)
	if error != nil {
		context.JSON(404, map[string]string{"msg": "Results Data not found!"})
		return
	}

	context.JSON(200, result)
}

func dbMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	_, error := strconv.Atoi(port)
	if error != nil {
		log.Fatalf("Invalid Port Value error : %v", error)
		panic("Invalid Port Value!")
	}

	postgresHost := os.Getenv("DB_HOST")
	postgresPort := os.Getenv("DB_PORT")
	_, error = strconv.Atoi(postgresPort)
	if error != nil {
		log.Fatalf("Invalid Port Value For DB : %v", error)
		panic("Invalid Port Value For DB!")
	}
	postgresUsername := os.Getenv("DB_USERNAME")
	postgresPassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUsername, postgresPassword, postgresHost, postgresPort, dbname)
	db, error := sql.Open("postgres", connStr)
	if error != nil {
		log.Fatalf("Failed on connect to postgres : %v", error)
		panic("Failed connect postgres!")
	}

	server := gin.Default()
	server.Use(dbMiddleware(db))

	server.GET("/job/:id", getJob)

	portString := fmt.Sprintf(":%s", port)
	server.Run(portString)
}

package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const REPO_URL = "https://raw.githubusercontent.com/quantum-plugins"

type GetJobById struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type AddPluginByName struct {
	Name string `uri:"name" binding:"required"`
}

func getJobData(jobId string, db *sql.DB) (sql.Result, error) {
	return db.Exec(`
		SELECT * 
		FROM results
		WHERE job_id=$1 
	`, jobId)
}

func getJob(context *gin.Context) {
	var job GetJobById
	err := context.ShouldBindUri(&job)
	if err != nil {
		context.JSON(400, map[string]string{"msg": err.Error()})
		return
	}

	db, ok := context.MustGet("db").(*sql.DB)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
	}

	result, err := getJobData(job.ID, db)
	if err != nil {
		context.JSON(404, map[string]string{"msg": "Results Data not found!"})
		return
	}

	context.JSON(200, result)
}

func getBackends(pluginName string) ([]string, error) {
	pipName := strings.Replace(pluginName, "-", "_", -1)
	fullBackendsListURL := fmt.Sprintf("%s/%s/refs/heads/main/%s/backends.txt", REPO_URL, pluginName, pipName)

	response, err := http.Get(fullBackendsListURL)
	if err != nil {
		return []string{}, err
	}
	defer response.Body.Close()

	backends, err := io.ReadAll(response.Body)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(backends), "\n")

	return lines, nil
}

func saveOnDB(backends *[]string, pluginName string, db *sql.DB) error {

	for _, backend := range *backends {
		_, err := db.Exec(`
			INSERT INTO backends(backend_name, plugin)
			VALUES($1, $2)
		`, backend, pluginName)

		if err != nil {
			return err
		}
	}

	return nil
}

func addPlugin(context *gin.Context) {
	var plugin AddPluginByName
	err := context.ShouldBindUri(&plugin)
	if err != nil {
		context.JSON(400, map[string]string{"msg": err.Error()})
	}

	db, ok := context.MustGet("db").(*sql.DB)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
	}

	pluginName := plugin.Name
	backends, err := getBackends(pluginName)

	if err != nil || len(backends) <= 0 {
		context.JSON(500, map[string]string{"msg": "Failed get backends!"})
	}

	err = saveOnDB(&backends, pluginName, db)
	if err != nil {
		context.JSON(500, map[string]string{"msg": "Failed on save data on DB!"})
	}

	context.JSON(201, map[string]string{"msg": "added plugin"})

}

func dbMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid Port Value error : %v", err)
		panic("Invalid Port Value!")
	}

	postgresHost := os.Getenv("DB_HOST")
	postgresPort := os.Getenv("DB_PORT")
	_, err = strconv.Atoi(postgresPort)
	if err != nil {
		log.Fatalf("Invalid Port Value For DB : %v", err)
		panic("Invalid Port Value For DB!")
	}
	postgresUsername := os.Getenv("DB_USERNAME")
	postgresPassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUsername, postgresPassword, postgresHost, postgresPort, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed on connect to postgres : %v", err)
		panic("Failed connect postgres!")
	}

	server := gin.Default()

	// check: https://stackoverflow.com/questions/34046194/how-to-pass-arguments-to-router-handlers-in-golang-using-gin-web-framework
	server.Use(dbMiddleware(db))

	server.GET("/job/:id", getJob)
	server.POST("/plugin/:name", addPlugin)

	portString := fmt.Sprintf(":%s", port)
	server.Run(portString)
}

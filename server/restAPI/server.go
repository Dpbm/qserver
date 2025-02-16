package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

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
	`, (jobId))
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

func addPlugin(context *gin.Context) {
	var plugin AddPluginByName
	err := context.ShouldBindUri(&plugin)
	if err != nil {
		context.JSON(400, map[string]string{"msg": err.Error()})
		return
	}

	db, ok := context.MustGet("db").(*sql.DB)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on Stablish database connection!"})
	}

	queueChannel, ok := context.MustGet("queueChannel").(*amqp.Channel)
	if !ok {
		context.JSON(500, map[string]string{"msg": "Failed on get queue channel!"})
	}

	// get all backends names with python script
	// go to the community official plugins, and get using curl, or whatever, the list
	// add to queue to install

}

func dbMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	}
}

func queueMiddleware(channel *amqp.Channel) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("queueChannel", channel)
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

	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	_, err = strconv.Atoi(rabbitmqPort)
	if err != nil {
		log.Fatalf("Invalid Port Value For RabbitMQ : %v", err)
		panic("Invalid Port Value For RabbitMQ!")
	}
	rabbitmqServerUrl := fmt.Sprintf("amqp://guest:guest@%s:%s", rabbitmqHost, rabbitmqPort)
	rabbitmqConnection, err := amqp.Dial(rabbitmqServerUrl)
	if err != nil {
		log.Fatalf("Failed on connect to rabbitmq : %v", err)
		panic("Failed rabbitmq connect")
	}
	defer rabbitmqConnection.Close()

	rabbitmqChannel, err := rabbitmqConnection.Channel()
	if err != nil {
		log.Fatalf("Failed on connect to rabbitmq channel : %v", err)
		panic("Failed rabbitmq channel")
	}
	defer rabbitmqChannel.Close()

	//----------------------------------------------------------------------------
	server := gin.Default()

	// check: https://stackoverflow.com/questions/34046194/how-to-pass-arguments-to-router-handlers-in-golang-using-gin-web-framework
	server.Use(dbMiddleware(db))
	server.Use(queueMiddleware(rabbitmqChannel))

	server.GET("/job/:id", getJob)
	server.POST("/plugin/:name", addPlugin)

	portString := fmt.Sprintf(":%s", port)
	server.Run(portString)
}

package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	constants "github.com/Dpbm/quantumRestAPI/constants"
	"github.com/Dpbm/quantumRestAPI/db"
	"github.com/Dpbm/quantumRestAPI/server"
	"github.com/stretchr/testify/assert"
)

func TestGetInvalidUUIDJob(t *testing.T) {
	os.Setenv("DB_PORT", "1") // to pass the por env check

	dbInstance := db.DB{}
	dbInstance.Connect(&db.Mock{})
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/job/nothing", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)

}

func TestGetUUIDNotFound(t *testing.T) {
	os.Setenv("DB_PORT", "1") // to pass the por env check

	dbInstance := db.DB{}
	dbInstance.Connect(&db.Mock{})
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/job/00000000-0000-0000-0000-000000000000", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)

	dbInstance.CloseConnection()

}

func TestFailDBConnection(t *testing.T) {
	var dbInstance *db.DB = nil
	server := server.SetupServer(dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/job/00000000-0000-0000-0000-000000000000", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)

}

func TestGetCorrectJob(t *testing.T) {
	os.Setenv("DB_PORT", "1") // to pass the por env check

	dbInstance := db.DB{}
	dbInstance.Connect(&db.Mock{})
	defer dbInstance.CloseConnection()

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	rows := sqlmock.NewRows([]string{"id", "job_id", "counts", "quasi_dist", "expval"}).
		AddRow(constants.TEST_JOB_ID, constants.TEST_JOB_ID, "{}", "{}", 10.3)

	mock.ExpectQuery("FROM results").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	log.Println(writer.Body.String())

	result := fmt.Sprintf(`{"id":"%s","job_id":"%s","counts":{},"quasi_dist":{},"expval":10.3}`, constants.TEST_JOB_ID, constants.TEST_JOB_ID)
	assert.Equal(t, result, writer.Body.String())
}

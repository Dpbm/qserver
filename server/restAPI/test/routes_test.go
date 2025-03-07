package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	constnats "github.com/Dpbm/quantumRestAPI/constants"
	"github.com/Dpbm/quantumRestAPI/db"
	"github.com/Dpbm/quantumRestAPI/server"
	dbDefinition "github.com/Dpbm/shared/db"
	"github.com/stretchr/testify/assert"
)

// TEST DELETE PLUGIN (WITHOUT ANY JOB ASSIGNED TO IT, on HISTORY)
// TEST DELETE WITH JOBS USING IT

const dbHost = ""
const dbPort = 0
const dbUsername = ""
const dbPassword = ""
const dbName = ""

// ------- ADD PLUGIN -------
func TestAddPluginFailedNoPluginName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestAddPluginFailedInvalidPlugin(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/invalid-plugin-test-it-should-not-work", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestAddPluginSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectExec("INSERT INTO backends").WithArgs("aer", "aer-plugin").WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/plugin/%s", constnats.TEST_PLUGIN), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)
}

// ------- ADD PLUGIN -------

func TestDeletePluginFailedNoPluginName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugin/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestDeletePluginFailedNotFoundInstalledPlugin(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugin/doesnt-exists", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestDeletePluginSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectExec("DELETE FROM backends").WithArgs("aer-plugin").WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugin/aer-plugin", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)
}

/*

func TestGetNoJobId(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)

}

func TestGetInvalidUUIDJob(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/nothing", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)

}

func TestGetUUIDNotFound(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/00000000-0000-0000-0000-000000000000", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)

	dbInstance.CloseConnection()

}

func TestFailDBConnectionGetJob(t *testing.T) {
	var dbInstance *db.DB = nil
	server := server.SetupServer(dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/00000000-0000-0000-0000-000000000000", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)

}

func TestGetCorrectJob(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
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
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	result := fmt.Sprintf(`{"id":"%s","job_id":"%s","counts":{},"quasi_dist":{},"expval":10.3}`, constants.TEST_JOB_ID, constants.TEST_JOB_ID)
	assert.Equal(t, result, writer.Body.String())
}

func TestAddPluginNoName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)

}

func TestGetInvalidPluginName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/invalid-name", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)

}

func TestFailDBConnectionAddPlugin(t *testing.T) {
	var dbInstance *db.DB = nil
	server := server.SetupServer(dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/plugin/%s", constants.TEST_PLUGIN), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestPotCorrectPlugin(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectExec("INSERT INTO backends").WithArgs("aer", "aer-plugin").WillReturnResult(sqlmock.NewResult(1, 1))

	server := server.SetupServer(&dbInstance)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/plugin/%s", constants.TEST_PLUGIN), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)

}
*/

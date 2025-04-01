package test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	constants "github.com/Dpbm/quantumRestAPI/constants"
	"github.com/Dpbm/quantumRestAPI/db"
	"github.com/Dpbm/quantumRestAPI/server"
	"github.com/Dpbm/quantumRestAPI/types"
	dbDefinition "github.com/Dpbm/shared/db"
	logger "github.com/Dpbm/shared/log"
	"github.com/stretchr/testify/assert"
)

const dbHost = ""
const dbPort = 0
const dbUsername = ""
const dbPassword = ""
const dbName = ""
const proxy = ""

// ADD A FAKE PLUGIN FOR THAT (use instead of aer)
const plugin_name = "aer-plugin"
const backend_name = "aer"

// ------- ADD PLUGIN -------
func TestAddPluginFailedNoPluginName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestAddPluginFailedInvalidPlugin(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugin/invalid-plugin-test-it-should-not-work", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestAddPluginSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectExec("INSERT INTO backends").WithArgs(backend_name, plugin_name).WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/plugin/%s", constants.TEST_PLUGIN), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)
}

func TestAddPluginNoRowsAffected(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectExec("INSERT INTO backends").WithArgs(backend_name, plugin_name).WillReturnResult(sqlmock.NewResult(1, 0))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/plugin/%s", constants.TEST_PLUGIN), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

// ------- DELETE PLUGIN -------

func TestDeletePluginFailedNoPluginName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugin/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestDeletePluginFailedNotFoundInstalledPlugin(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	const invalidName string = "doesnt-exits"
	row := mock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT").WithArgs(invalidName).WillReturnRows(row)
	mock.ExpectExec("DELETE FROM backends").WithArgs(invalidName).WillReturnResult(sqlmock.NewResult(1, 0))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/plugin/%s", invalidName), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestDeletePluginSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	row := mock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT").WithArgs(plugin_name).WillReturnRows(row)
	mock.ExpectExec("DELETE FROM backends").WithArgs(plugin_name).WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/plugin/%s", plugin_name), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)
}

func TestDeletePluginFailedRunningJobs(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	row := mock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT").WithArgs(plugin_name).WillReturnRows(row)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/plugin/%s", plugin_name), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

// ------- GET BACKEND -------

func TestGetBackendFailedNoBackendName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/backend/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetBackendFailedNoBackendWithThisName(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	const invalidName string = "invalid-backend"
	mock.ExpectQuery("FROM backends").WithArgs(invalidName).WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/backend/%s", invalidName), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetBackendSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	rows := mock.NewRows([]string{"backend_name", "id", "pointer", "plugin"}).AddRow(constants.TEST_BACKEND, "1", 1, constants.TEST_PLUGIN)
	mock.ExpectQuery("FROM backends").WithArgs(constants.TEST_BACKEND).WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/backend/%s", constants.TEST_BACKEND), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body types.BackendData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, body.ID, "1")
	assert.Equal(t, body.Name, constants.TEST_BACKEND)
	assert.Equal(t, body.Plugin, constants.TEST_PLUGIN)
	assert.Equal(t, body.Pointer, uint64(1))
}

// ------- GET BACKENDS -------

func TestNoBackends(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("FROM backends").WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/backends/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.BackendData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, len(body), 0)
}

func TestOneBackend(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	values := [][]driver.Value{
		{
			constants.TEST_BACKEND, "1", 1, constants.TEST_PLUGIN,
		},
	}
	rows := mock.NewRows([]string{"backend_name", "id", "pointer", "plugin"}).AddRows(values...)
	mock.ExpectQuery("FROM backends").WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/backends/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.BackendData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, len(body), 1)
	assert.Equal(t, body[0].ID, "1")
	assert.Equal(t, body[0].Name, constants.TEST_BACKEND)
	assert.Equal(t, body[0].Plugin, constants.TEST_PLUGIN)
	assert.Equal(t, body[0].Pointer, uint64(1))
}

// ------- GET JOB -------

func TestGetJobWithoutPassingJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetJobUsingInvalidJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/invalid-id", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)
}

func TestGetJobIDNotFound(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetJobSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now().UTC().Local()

	rows := mock.NewRows([]string{
		"id",
		"pointer",
		"target_simulator",
		"qasm",
		"status",
		"submission_date",
		"start_time",
		"finish_time",
		"metadata",
		"result_types",
		"results",
	}).AddRow(constants.TEST_JOB_ID, 1, constants.TEST_BACKEND, "nothing", "pending", now, now, now, "{}", "{}", "{}")
	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	expectedMetadata := map[string]any{}
	expectedResultTypes := types.JobResultTypes{ID: "", JobId: "", Counts: false, QuasiDist: false, Expval: false}
	expectedResults := types.JobResultData{ID: "", JobId: "", Counts: map[string]float64(nil), QuasiDist: map[int64]float64(nil), Expval: []float64(nil)}

	var data types.JobData
	err := json.NewDecoder(writer.Result().Body).Decode(&data)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, data.ID, constants.TEST_JOB_ID)
	assert.Equal(t, data.FinishTime.Time.String(), now.String())
	assert.Equal(t, data.Metadata, expectedMetadata)
	assert.Equal(t, data.Pointer, uint64(1))
	assert.Equal(t, data.Qasm, "nothing")
	assert.Equal(t, data.ResultTypes, expectedResultTypes)
	assert.Equal(t, data.Results, expectedResults)
	assert.Equal(t, data.StartTime.Time.String(), now.String())
	assert.Equal(t, data.Status, "pending")
	assert.Equal(t, data.SubmissionDate.String(), now.String())
	assert.Equal(t, data.TargetSimulator, constants.TEST_BACKEND)
}

// ------- GET JOB RESULT -------

func TestGetJobResultWithoutPassingJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/result/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetJobResultUsingInvalidJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/result/invalid-id", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)
}

func TestGetJobResultIDNotFound(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("FROM results").WithArgs(constants.TEST_JOB_ID).WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/job/result/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestGetJobResultSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	rows := mock.NewRows([]string{
		"id",
		"job_id",
		"counts",
		"quasi_dist",
		"expval",
	}).AddRow("1", constants.TEST_JOB_ID, "{}", "{}", "[]")
	mock.ExpectQuery("FROM results").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/job/result/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var data types.JobResultData
	err := json.NewDecoder(writer.Result().Body).Decode(&data)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, data.ID, "1")
	assert.Equal(t, data.JobId, constants.TEST_JOB_ID)
	assert.Equal(t, data.Counts, map[string]float64{})
	assert.Equal(t, data.QuasiDist, map[int64]float64{})
	assert.Equal(t, data.Expval, []float64{})
}

// ------- CANCEL JOB -------

func TestCancelJobWithoutPassingJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/job/cancel/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestCancelJobWithInvalidID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/job/cancel/dadadas", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)
}

func TestCancelJobIDNotFound(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/job/cancel/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestCancelJobSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now()
	rows := mock.NewRows([]string{
		"id",
		"pointer",
		"target_simulator",
		"qasm",
		"status",
		"submission_date",
		"start_time",
		"finish_time",
		"metadata",
		"result_types",
		"results",
	}).AddRow(constants.TEST_JOB_ID, 1, constants.TEST_BACKEND, "nothing", "pending", now, now, now, "{}", "{}", "{}")
	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	mock.ExpectExec("UPDATE jobs").WithArgs(constants.TEST_JOB_ID).WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/job/cancel/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)
}

func TestCancelErrorJobRunning(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now()
	rows := mock.NewRows([]string{
		"id",
		"pointer",
		"target_simulator",
		"qasm",
		"status",
		"submission_date",
		"start_time",
		"finish_time",
		"metadata",
		"result_types",
		"results",
	}).AddRow(constants.TEST_JOB_ID, 1, constants.TEST_BACKEND, "nothing", "running", now, now, now, "{}", "{}", "{}")
	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/job/cancel/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestCancelErrorJobFinished(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now()
	rows := mock.NewRows([]string{
		"id",
		"pointer",
		"target_simulator",
		"qasm",
		"status",
		"submission_date",
		"start_time",
		"finish_time",
		"metadata",
		"result_types",
		"results",
	}).AddRow(constants.TEST_JOB_ID, 1, constants.TEST_BACKEND, "nothing", "finished", now, now, now, "{}", "{}", "{}")
	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/job/cancel/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestCancelErrorNoRowsAffected(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now()
	rows := mock.NewRows([]string{
		"id",
		"pointer",
		"target_simulator",
		"qasm",
		"status",
		"submission_date",
		"start_time",
		"finish_time",
		"metadata",
		"result_types",
		"results",
	}).AddRow(constants.TEST_JOB_ID, 1, constants.TEST_BACKEND, "nothing", "pending", now, now, now, "{}", "{}", "{}")
	mock.ExpectQuery("FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnRows(rows)
	mock.ExpectExec("UPDATE jobs").WithArgs(constants.TEST_JOB_ID).WillReturnResult(sqlmock.NewResult(1, 0))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/job/cancel/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

// ------- DELETE JOB -------

func TestDeleteJobWithoutPassingJobID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/job/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 404, writer.Code)
}

func TestDeleteJobPassingInvalidID(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/job/invalid-id", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 400, writer.Code)
}

func TestDeleteJobIDNotFound(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestDeleteSuccess(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	row := mock.NewRows([]string{"status"}).AddRow("finished")
	mock.ExpectQuery("SELECT status").WithArgs(constants.TEST_JOB_ID).WillReturnRows(row)
	mock.ExpectExec("DELETE FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnResult(sqlmock.NewResult(1, 1))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)
}

func TestDeleteFailedJobRunning(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	row := mock.NewRows([]string{"status"}).AddRow("running")
	mock.ExpectQuery("SELECT status FROM").WithArgs(constants.TEST_JOB_ID).WillReturnRows(row)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestDeleteNoReturnRowsOnGet(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("SELECT status").WithArgs(constants.TEST_JOB_ID).WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

func TestDeleteNoAffectedRows(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)
	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	row := mock.NewRows([]string{"status"}).AddRow("finished")
	mock.ExpectQuery("SELECT status").WithArgs(constants.TEST_JOB_ID).WillReturnRows(row)
	mock.ExpectExec("DELETE FROM jobs").WithArgs(constants.TEST_JOB_ID).WillReturnResult(sqlmock.NewResult(1, 0))

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/job/%s", constants.TEST_JOB_ID), nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 500, writer.Code)
}

// ------- GET JOBS -------

func TestNoJobs(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("FROM backends").WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/jobs/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.BackendData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, len(body), 0)
}

func TestOneJob(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now().UTC().Local()
	values := [][]driver.Value{
		{
			"1",
			1,
			constants.TEST_BACKEND,
			"AAAA",
			"pending",
			now,
			now,
			now,
			"{}",
			"{}",
			"{}",
		},
	}
	rows := mock.NewRows([]string{"id", "pointer", "target_simulator", "qasm", "status", "submission_date", "start_time", "finish_time", "metadata", "result_types", "results"}).AddRows(values...)
	mock.ExpectQuery("FROM jobs").WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/jobs/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.JobData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	expectedMetadata := map[string]any{}
	expectedResultTypes := types.JobResultTypes{ID: "", JobId: "", Counts: false, QuasiDist: false, Expval: false}
	expectedResults := types.JobResultData{ID: "", JobId: "", Counts: map[string]float64(nil), QuasiDist: map[int64]float64(nil), Expval: []float64(nil)}

	assert.Equal(t, len(body), 1)
	assert.Equal(t, body[0].ID, "1")
	assert.Equal(t, body[0].Pointer, uint64(1))
	assert.Equal(t, body[0].TargetSimulator, constants.TEST_BACKEND)
	assert.Equal(t, body[0].Qasm, "AAAA")
	assert.Equal(t, body[0].Status, "pending")
	assert.Equal(t, body[0].SubmissionDate.String(), now.String())
	assert.Equal(t, body[0].StartTime.Time.String(), now.String())
	assert.Equal(t, body[0].FinishTime.Time.String(), now.String())
	assert.Equal(t, body[0].Metadata, expectedMetadata)
	assert.Equal(t, body[0].ResultTypes, expectedResultTypes)
	assert.Equal(t, body[0].Results, expectedResults)
}

// ------- GET HISTORY -------

func TestNoHistory(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	mock.ExpectQuery("FROM history").WillReturnError(sql.ErrNoRows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/history/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.BackendData
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	assert.Equal(t, len(body), 0)
}

func TestOneHistoryJob(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	mock, ok := dbInstance.Extra.(sqlmock.Sqlmock)
	if !ok {
		t.Fatal("Failed on parse mock")
	}

	now := time.Now().UTC().Local()
	values := [][]driver.Value{
		{
			"1",
			constants.TEST_JOB_ID,
			constants.TEST_BACKEND,
			"AAAA",
			"pending",
			now,
			now,
			now,
			"{}",
			"{}",
			"{}",
		},
	}
	rows := mock.NewRows([]string{"id", "job_id", "target_simulator", "qasm", "status", "submission_date", "start_time", "finish_time", "metadata", "result_types", "results"}).AddRows(values...)
	mock.ExpectQuery("FROM history").WillReturnRows(rows)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/history/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	var body []types.Historydata
	err := json.NewDecoder(writer.Result().Body).Decode(&body)

	if err != nil {
		logger.LogError(errors.New("decoding Failed"))
	}

	expectedMetadata := map[string]any{}
	expectedResultTypes := types.JobResultTypes{ID: "", JobId: "", Counts: false, QuasiDist: false, Expval: false}
	expectedResults := types.JobResultData{ID: "", JobId: "", Counts: map[string]float64(nil), QuasiDist: map[int64]float64(nil), Expval: []float64(nil)}

	assert.Equal(t, len(body), 1)
	assert.Equal(t, body[0].ID, uint64(1))
	assert.Equal(t, body[0].JobId, constants.TEST_JOB_ID)
	assert.Equal(t, body[0].TargetSimulator, constants.TEST_BACKEND)
	assert.Equal(t, body[0].Qasm, "AAAA")
	assert.Equal(t, body[0].Status, "pending")
	assert.Equal(t, body[0].SubmissionDate.String(), now.String())
	assert.Equal(t, body[0].StartTime.Time.String(), now.String())
	assert.Equal(t, body[0].FinishTime.Time.String(), now.String())
	assert.Equal(t, body[0].Metadata, expectedMetadata)
	assert.Equal(t, body[0].ResultTypes, expectedResultTypes)
	assert.Equal(t, body[0].Results, expectedResults)
}

// ------- GET HEALTH -------

func TestHealth(t *testing.T) {
	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Mock{}, dbHost, dbPort, dbUsername, dbPassword, dbName)
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance, proxy)

	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health/", nil)

	server.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

}

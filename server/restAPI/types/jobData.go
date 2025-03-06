package types

import (
	"database/sql"
	"time"
)

type Metada = map[string]any

type JobResultData struct {
	ID        string             `json:"id"`
	JobId     string             `json:"job_id"`
	Counts    map[string]float32 `json:"counts"`
	QuasiDist map[int32]float32  `json:"quasi_dist"`
	Expval    []float32          `json:"expval"`
}

type JobResultTypes struct {
	ID        string `json:"id"`
	JobId     string `json:"job_id"`
	Counts    bool   `json:"counts"`
	QuasiDist bool   `json:"quasi_dist"`
	Expval    bool   `json:"expval"`
}

type JobData struct {
	ID              string         `json:"id"`
	Pointer         uint32         `json:"pointer"`
	TargetSimulator string         `json:"target_simulator"`
	Qasm            string         `json:"qasm"`
	Status          string         `json:"status"`
	SubmissionDate  time.Time      `json:"submission_date"`
	StartTime       sql.NullTime   `json:"start_time"`
	FinishTime      sql.NullTime   `json:"finish_time"`
	Metadata        Metada         `json:"metadata"`
	ResultTypes     JobResultTypes `json:"result_types"`
	Results         JobResultData  `json:"results"`
}

type Historydata struct {
	ID              uint32         `json:"id"`
	JobId           string         `json:"job_id"`
	TargetSimulator string         `json:"target_simulator"`
	Qasm            string         `json:"qasm"`
	Status          string         `json:"status"`
	SubmissionDate  time.Time      `json:"submission_date"`
	StartTime       sql.NullTime   `json:"start_time"`
	FinishTime      sql.NullTime   `json:"finish_time"`
	Metadata        Metada         `json:"metadata"`
	ResultTypes     JobResultTypes `json:"result_types"`
	Results         JobResultData  `json:"results"`
}

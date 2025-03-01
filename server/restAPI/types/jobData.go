package types

import "time"

type JobResultData struct {
	ID        string             `json:"id"`
	JobId     string             `json:"job_id"`
	Counts    map[string]float32 `json:"counts"`
	QuasiDist map[int32]float32  `json:"quasi_dist"`
	Expval    []float32          `json:"expval"`
}

type JobData struct {
	ID              string            `json:"id"`
	Order           uint32            `json:"order"`
	TargetSimulator string            `json:"target_simulator"`
	Qasm            string            `json:"qasm"`
	Status          string            `json:"status"`
	SubmissionDate  time.Time         `json:"submission_date"`
	StartTime       time.Time         `json:"start_time"`
	FinishTime      time.Time         `json:"finish_time"`
	Metadata        map[any]any       `json:"metadata"`
	ResultTypes     map[string]string `json:"result_types"`
	Results         map[string]any    `json:"results"`
}

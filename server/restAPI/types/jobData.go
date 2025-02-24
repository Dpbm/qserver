package types

type JobData struct {
	ID        string             `json:"id"`
	JobId     string             `json:"job_id"`
	Counts    map[string]int32   `json:"counts"`
	QuasiDist map[string]float32 `json:"quasi_dist"`
	Expval    float32            `json:"expval"`
}

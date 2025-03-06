package data

import (
	"encoding/json"
	"errors"

	jobsServerProto "github.com/Dpbm/jobsServer/proto"
)

func CheckData(data *jobsServerProto.JobProperties) error {
	if data == nil {
		return errors.New("invalid job properties")
	}

	if !data.ResultTypeCounts && !data.ResultTypeExpVal && !data.ResultTypeQuasiDist {
		return errors.New("you must select at least one type of result")
	}

	if len(data.TargetSimulator) <= 0 {
		return errors.New("you must the target simulator to run your job")
	}

	if len(data.Metadata) <= 0 {
		return errors.New("invalid metadata")
	}

	_, err := json.Marshal(data.Metadata)

	return err
}

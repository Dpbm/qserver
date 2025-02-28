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

	// refer to: https://stackoverflow.com/questions/43770273/json-unmarshalling-without-struct
	// and: https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json
	var metadata map[string]interface{}
	if data.Metadata != nil {
		err := json.Unmarshal([]byte(*data.Metadata), &metadata)

		if len(*data.Metadata) <= 0 || err != nil {
			return errors.New("invalid metadata")
		}
	}

	return nil
}

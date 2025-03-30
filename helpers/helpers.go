package helpers

import (
	"bugtigexa.giantfilesuploader.com/model"
	"encoding/json"
	"fmt"
)

const (
	DATA string = "data"
)

func Decode(b []byte) (interface{}, error) {
	objs := map[string]interface{}{
		DATA: &model.Data{},
	}
	for _, v := range objs {
		if err := json.Unmarshal(b, v); err == nil {
			return v, nil
		}
	}
	return nil, fmt.Errorf("unable to unmarshal json")
}

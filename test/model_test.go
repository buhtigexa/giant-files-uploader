package test

import (
	"bugtigexa.giantfilesuploader.com/model"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalModel(t *testing.T) {
	data := &model.Data{
		Id:    "id",
		Part:  0,
		Value: []byte("value"),
		Time:  time.Now().UTC(),
	}

	by, err := json.Marshal(data)
	assert.Nil(t, err)

	var udata model.Data
	err = json.Unmarshal(by, &udata)
	assert.Nil(t, err)
	assert.Equal(t, data.Id, udata.Id)
	assert.Equal(t, data.Part, udata.Part)
	assert.Equal(t, data.Value, udata.Value)
	assert.Equal(t, data.Time.Unix(), udata.Time.Unix())
}

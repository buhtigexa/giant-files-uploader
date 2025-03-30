package test

import (
	"bugtigexa.giantfilesuploader.com/cmd"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalModel(t *testing.T) {
	data := &cmd.Data{
		FileName: "id",
		Part:     0,
		Value:    []byte("value"),
		Time:     time.Now().UTC(),
	}

	by, err := json.Marshal(data)
	assert.Nil(t, err)

	var udata cmd.Data
	err = json.Unmarshal(by, &udata)
	assert.Nil(t, err)
	assert.Equal(t, data.FileName, udata.FileName)
	assert.Equal(t, data.Part, udata.Part)
	assert.Equal(t, data.Value, udata.Value)
	assert.Equal(t, data.Time.Unix(), udata.Time.Unix())
}

package model

import (
	"bugtigexa.giantfilesuploader.com/helpers"
	"encoding/json"
	"fmt"
)

type FileManager struct {
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (m *FileManager) Store(b []byte) (int, error) {
	var data Data
	if err := json.Unmarshal(b, &data); err != nil {
		return 0, err
	}

	obj, err := helpers.Decode(b)
	if err != nil {
		return 0, err
	}

	switch v := obj.(type) {
	case Data:
		data = v
	default:
		return 0, fmt.Errorf("unable to unmarshal json")
	}
	if _, err := data.Save(); err != nil {
		return 0, err
	}

	return 0, nil
}

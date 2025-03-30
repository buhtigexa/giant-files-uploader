package cmd

import (
	"encoding/json"
	"fmt"
)

type FileManager struct {
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (m *FileManager) Store(b []byte) (int, error) {
	obj, err := Decode(b)
	if err != nil {
		return 0, err
	}

	switch v := obj.(type) {
	case *Data:
		n, err := v.Save()
		if err != nil {
			return 0, err
		}
		return n, err
	default:
		return 0, fmt.Errorf("unable to unmarshal json")
	}

	return 0, nil
}

func Decode(b []byte) (obj interface{}, err error) {
	json.Unmarshal(b, &obj)
	return obj, nil
}

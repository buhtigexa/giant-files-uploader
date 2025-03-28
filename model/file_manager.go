package model

import (
	"encoding/json"
	"os"
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

	dir := GetDirname(data)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return 0, err
		}
	}

	fname := GetFileName(data)
	_, err = os.Stat(fname)
	if !os.IsNotExist(err) {
		return 0, nil
	}
	f, err := os.Create(fname)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err := f.Write(b)
	if err != nil {
		return n, err
	}
	return n, nil
}

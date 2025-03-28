package model

import (
	"fmt"
)

func GetFileName(data Data) string {
	return fmt.Sprintf("%s/%s-%d-%d", GetDirname(data), data.Id, data.Part, data.Time.UTC().Unix())
}

func GetDirname(data Data) string {
	return fmt.Sprintf("%s", data.Name)
}

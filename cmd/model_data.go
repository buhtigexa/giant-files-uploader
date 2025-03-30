package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Data struct {
	FileName string    `json:"id"`
	Part     int       `json:"part"`
	Value    []byte    `json:"value"`
	Time     time.Time `json:"time"`
	Total    float32   `json:"total,omitempty"`
}

func (data *Data) MarshalJSON() ([]byte, error) {
	uxtime := data.Time.UTC().Unix()
	aux := struct {
		Id    string  `json:"id"`
		Part  int     `json:"part"`
		Value []byte  `json:"value"`
		Time  int64   `json:"time"`
		Total float32 `json:"total,omitempty"`
	}{
		Id:    data.FileName,
		Part:  data.Part,
		Value: data.Value,
		Time:  uxtime,
		Total: data.Total,
	}
	by, err := json.Marshal(aux)
	return by, err
}

func (data *Data) UnmarshalJSON(b []byte) error {
	var aux struct {
		Id    string  `json:"id"`
		Part  int     `json:"part"`
		Value []byte  `json:"value"`
		Time  int64   `json:"time"`
		Total float32 `json:"total,omitempty"`
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	data.FileName = aux.Id
	data.Part = aux.Part
	data.Value = aux.Value
	data.Time = time.Unix(aux.Time, 0).UTC()
	data.Total = aux.Total
	return nil
}

func (d *Data) Save() (int, error) {
	dir := d.GetDirname()
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return 0, err
		}
	}

	fname := d.GetFileName()
	_, err = os.Stat(fname)
	if !os.IsNotExist(err) {
		return 0, nil
	}
	f, err := os.Create(fname)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err := f.Write(d.Value)
	return n, err
}

func (d *Data) toDataEntity() DataEntity {
	return DataEntity{
		FileName: d.FileName,
		Total:    d.Total,
		Time:     d.Time.UTC().Unix(),
	}
}

func (d Data) GetFileName() string {
	return fmt.Sprintf("%s/%s-%d-%d", d.GetDirname(), d.FileName, d.Part, d.Time.UTC().Unix())
}

func (d Data) GetDirname() string {
	return fmt.Sprintf("%s", d.FileName)
}

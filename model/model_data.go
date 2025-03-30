package model

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Data struct {
	Id    string    `json:"id"`
	Part  int       `json:"part"`
	Value []byte    `json:"value"`
	Time  time.Time `json:"time"`
	Total float32   `json:"total"`
}

func (data *Data) MarshalJSON() ([]byte, error) {
	uxtime := data.Time.UTC().Unix()
	aux := struct {
		Id    string  `json:"id"`
		Part  int     `json:"part"`
		Value []byte  `json:"value"`
		Time  int64   `json:"time"`
		Total float32 `json:"total"`
	}{
		Id:    data.Id,
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
		Total float32 `json:"total"`
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	data.Id = aux.Id
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

func (d Data) GetFileName() string {
	return fmt.Sprintf("%s/%s-%d-%d", d.GetDirname(), d.Id, d.Part, d.Time.UTC().Unix())
}

func (d Data) GetDirname() string {
	return fmt.Sprintf("%s", d.Id)
}

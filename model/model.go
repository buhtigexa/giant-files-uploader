package model

import (
	"encoding/json"
	"time"
)

type Data struct {
	Id    string    `json:"id"`
	Part  int       `json:"part"`
	Name  string    `json:"name"`
	Value []byte    `json:"value"`
	Time  time.Time `json:"time"`
}

func (data *Data) MarshalJSON() ([]byte, error) {
	uxtime := data.Time.UTC().Unix()
	aux := struct {
		Id    string `json:"id"`
		Part  int    `json:"part"`
		Name  string `json:"name"`
		Value []byte `json:"value"`
		Time  int64  `json:"time"`
	}{
		Id:    data.Id,
		Part:  data.Part,
		Name:  data.Name,
		Value: data.Value,
		Time:  uxtime,
	}
	by, err := json.Marshal(aux)
	return by, err
}

func (data *Data) UnmarshalJSON(b []byte) error {
	var aux struct {
		Id    string `json:"id"`
		Part  int    `json:"part"`
		Name  string `json:"name"`
		Value []byte `json:"value"`
		Time  int64  `json:"time"`
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	data.Id = aux.Id
	data.Part = aux.Part
	data.Name = aux.Name
	data.Value = aux.Value
	data.Time = time.Unix(aux.Time, 0).UTC()
	return nil
}

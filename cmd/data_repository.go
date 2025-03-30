package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type DataEntity struct {
	ID       int
	FileName string
	Total    float32
	Time     int64
}

func (s *StreamServer) Save(data DataEntity) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var insertedID int64

	err := s.db.sqlDB.QueryRowContext(
		ctx,
		"INSERT INTO data (filename, total, time) VALUES ($1, $2, $3) RETURNING id",
		data.FileName, data.Total, data.Time,
	).Scan(&insertedID)

	if err != nil {
		log.Printf("Error while saving data: %v", err)
		return 0, err
	}

	return insertedID, nil
}

func (s *StreamServer) Load(id int64) (DataEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data DataEntity
	err := s.db.sqlDB.QueryRowContext(ctx, "SELECT id,fileName,total,time FROM data WHERE id=$1", id).Scan(&data.ID, &data.FileName, &data.Total, &data.Time)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return DataEntity{}, err
		}
		return DataEntity{}, err
	}
	return data, nil
}

func (m *StreamServer) GetAll(filters Filters) ([]DataEntity, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, filename, total, time
		FROM data
		ORDER BY %s
		LIMIT $1 OFFSET $2`, filters.SortColumn())

	args := []interface{}{filters.Limit(), filters.Offset()}

	rows, err := m.db.sqlDB.Query(query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	data := []DataEntity{}

	for rows.Next() {
		var d DataEntity
		err := rows.Scan(&totalRecords, &d.ID, &d.FileName, &d.Total, &d.Time)
		if err != nil {
			return nil, Metadata{}, err
		}
		data = append(data, d)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return data, metadata, nil
}

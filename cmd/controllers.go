package cmd

import (
	"encoding/json"
	"net/http"
)

func (s *StreamServer) GetUploads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	v := NewValidator()
	qs := r.URL.Query()
	page := readInt(qs, "page", 1, v)
	pageSize := readInt(qs, "pageSize", 10, v)
	sort := readString(qs, "sort", "id")

	if ok := v.Valid(); !ok {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "invalid query parameters",
			"fields": v.Errors,
		},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	sortSafeList := []string{"id", "filename", "total", "time", "-id", "-filename", "-total", "-time"}
	filters := Filters{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sort,
		SortSafeList: sortSafeList,
	}

	ds, metadata, err := s.GetAll(filters)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"data": ds,
		"meta": metadata,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

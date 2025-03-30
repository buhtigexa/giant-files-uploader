package cmd

import (
	"net/http"
)

func (s *StreamServer) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/uploads", http.HandlerFunc(s.GetUploads))
	return mux
}

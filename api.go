package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type ApiService struct {
	svc Service
}

func NewApiService(svc Service) *ApiService {
	return &ApiService{
		svc: svc,
	}
}

func (s *ApiService) start(listenAddr string) error {
	http.HandleFunc("/", s.handleGetMessage)
	return http.ListenAndServe(listenAddr, nil)
}

func (s *ApiService) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	fact, err := s.svc.GetMessage(context.Background())
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, fact)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

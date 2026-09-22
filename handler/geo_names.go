package handler

import (
	"encoding/json"
	"hsbc/services"
	"log"
	"net/http"
)

type response struct {
	Count int `json:"count"`
}

type GeoNamesHandler struct {
	api *services.GeoNamesApi
}

func (h GeoNamesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	letter := r.URL.Query().Get("letter")

	if len(letter) != 1 {
		http.Error(w, "expected exactly one character", http.StatusBadRequest)
		return
	}

	char := letter[0]
	if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
		http.Error(w, "expected a letter", http.StatusBadRequest)
		return
	}

	count, err := h.api.CountCitiesByLetter(letter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := response{
		Count: count,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func NewHandler(api *services.GeoNamesApi) http.Handler {
	return GeoNamesHandler{api: api}
}

package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func GetActor(r *http.Request) string {
	//supposed to be from session
	vars := mux.Vars(r)
	name, _ := vars["actor"]
	return name
}

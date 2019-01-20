package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
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

func GetBoolParam(r *http.Request, param string) bool {
	var includeRelated bool
	includeRelated, err := strconv.ParseBool(r.FormValue(param))
	if err != nil {
		includeRelated = false
	}
	return includeRelated
}

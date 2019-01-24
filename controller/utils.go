package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 10
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

func getActor(r *http.Request) string {
	//supposed to be from session
	vars := mux.Vars(r)
	name, _ := vars["actor"]
	return name
}

func getBoolParam(r *http.Request, param string) bool {
	var includeRelated bool
	includeRelated, err := strconv.ParseBool(r.FormValue(param))
	if err != nil {
		includeRelated = false
	}
	return includeRelated
}

func getPagingValues(r *http.Request) (int, int, error) {
	values := r.Form
	page := 1
	limit := DefaultLimit

	pageVal := values.Get("page")
	if len(pageVal) > 0 {
		page64, err := strconv.Atoi(pageVal)
		page = int(page64)
		if err != nil {
			return 0, 0, err
		}
	}

	limitVal := values.Get("limit")
	if len(limitVal) > 0 {
		limit64, err := strconv.ParseInt(limitVal, 10, 32)
		limit = int(limit64)
		if err != nil {
			return 0, 0, err
		}
	}

	return page, limit, nil
}

func getNextUrl(r *http.Request, page int, limit int) string {
	url := *r.URL
	values := url.Query()
	values.Set("page", strconv.Itoa(page+1))
	values.Set("limit", strconv.Itoa(limit))
	return fmt.Sprintf("http://%s%s?%s", r.Host, url.Path, values.Encode())
}

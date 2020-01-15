package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
)

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projects := fetchProjects()
	json.NewEncoder(w).Encode(projects)
}

func parseIsoDate(date string) time.Time {
	t, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return time.Now()
	}
	return t
}

func getQuery(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys) == 0 {
		return ""
	}
	return keys[0]
}

func getCommits(w http.ResponseWriter, r *http.Request) {
	fromQ := getQuery(r, "from")
	toQ := getQuery(r, "to")
	messageQ := getQuery(r, "message")
	from := parseIsoDate(fromQ)
	to := parseIsoDate(toQ)
	message, _ := regexp.Compile(messageQ)

	w.Header().Set("Content-Type", "application/json")
	commits := fetchCommits(&FetchCommitsOptions{from: &from, to: &to, withStats: true, messageRegex: message})
	json.NewEncoder(w).Encode(commits)
}

func foo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sometihgn")
	x, _ := r.URL.Query()["x"]
	log.Println(x)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(x)
}

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/projects", getProjects).Methods("GET")
	router.HandleFunc("/commits", getCommits).Methods("GET")
	router.HandleFunc("/bar", foo).Methods("GET")

	return router
}

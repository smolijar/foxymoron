package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	readConfig()
	port := 8000
	log.Printf("Startig server on port %v", port)
	http.ListenAndServe(":"+strconv.Itoa(port), createRouter())
}

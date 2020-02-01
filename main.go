package main

import (
	"log"
	"strconv"
)

func main() {
	readConfig()
	port := config.port
	log.Printf("Startig server on port %v", port)
	createRouter().Run(":" + strconv.Itoa(port))
}

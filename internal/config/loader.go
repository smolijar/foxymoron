package config

import (
	"flag"
	"os"
	"strconv"
)

var Config struct {
	Port int
}

func LoadConfig() {
	defaultPort := 8000
	envPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		defaultPort = envPort
	}
	port := flag.Int("port", defaultPort, "Local server port")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}
	Config.Port = *port
}

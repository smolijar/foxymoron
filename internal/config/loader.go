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
	// COOL: having built in param parsing is just cool
	port := flag.Int("port", defaultPort, "Local server port")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help == true {
		// COOL: print flag usage
		flag.PrintDefaults()
		os.Exit(0)
	}
	Config.Port = *port
}

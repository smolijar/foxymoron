package config

import (
	"flag"
	"os"
)

var Config struct {
	Port int
}

func LoadConfig() {
	// COOL: having built in param parsing is just cool
	port := flag.Int("port", 8000, "Local server port")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help == true {
		// COOL: print flag usage
		flag.PrintDefaults()
		os.Exit(0)
	}
	Config.Port = *port
}

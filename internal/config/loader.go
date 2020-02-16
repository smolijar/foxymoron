package config

import (
	"flag"
	"fmt"
	"os"
)

var Config struct {
	Token string
	Url   string
	Port  int
}

func LoadConfig() {
	// COOL: having built in param parsing is just cool
	token := flag.String("token", "", "GitLab API token")
	url := flag.String("url", "http://gitlab.com", "GitLab URL")
	port := flag.Int("port", 8000, "Local server port")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help == true {
		// COOL: print flag usage
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *token == "" {
		fmt.Println("Missing token")
		flag.PrintDefaults()
		os.Exit(1)
	}
	Config.Token = *token
	Config.Url = *url
	Config.Port = *port
}

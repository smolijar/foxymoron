package main

import (
	"github.com/grissius/foxymoron/internal/api"
	"github.com/grissius/foxymoron/internal/config"
)

func init() {
	config.LoadConfig()
}

func main() {
	api.RunAt(config.Config.Port)
}

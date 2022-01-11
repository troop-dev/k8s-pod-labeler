package main

import (
	"log"

	"github.com/troop-dev/k8s-pod-labeler/app/server"
)

func main() {
	cfg, err := server.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	server.Run(cfg)
}

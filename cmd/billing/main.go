package main

import (
	"github.com/iseroukhov/brave-new-billing/pkg/server"
	"log"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	s := server.New(cfg)
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}

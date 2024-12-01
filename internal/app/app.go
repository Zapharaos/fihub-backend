package app

import (
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"log"
)

func Init() {

	// Load the .env file
	err := env.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO : Services start
}

func Stop() {
	//	TODO : Services stop
}

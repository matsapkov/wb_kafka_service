package main

import (
	"github.com/matsapkov/wb_kafka_service/internal/server"
	"log"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	s.StartServer()
}

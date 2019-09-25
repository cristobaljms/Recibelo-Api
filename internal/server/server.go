package server

import (
	"log"
	"net/http"
	"recibe_me/configs"
)

// RunServer runs HTTP gateway
func RunServer() error {
	router := NewRouter()
	server := http.ListenAndServe(configs.DefaultServerConfig.Port, router)
	log.Fatal(server)
	return nil
}

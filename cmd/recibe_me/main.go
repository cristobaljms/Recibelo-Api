package main

import (
	"fmt"
	"log"
	"os"
	"recibe_me/internal/server"
)

func main() {
	if err := server.RunServer(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%v\n", err)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

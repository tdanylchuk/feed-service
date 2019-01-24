package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	server := CreateApp()
	defer server.Close()
	server.StartServer()
}

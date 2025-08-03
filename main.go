package main

import (
	"fmt"
)

func main() {
	router := NewRouter()
	router.HandlerFunc("/test", func(w ResponseWriter, r *Request) {
		fmt.Println("Test./........")
	})
	server := NewServer()
	fmt.Println("Starting to listen at 8080...")
	if err := server.StartListening("localhost:8080", router); err != nil {
		panic(err)
	}
}

// First create a connection to the client using server
// Parse its data.
//
// create a server that listens to specific given port for incoming request
// passes this request to handler.
//

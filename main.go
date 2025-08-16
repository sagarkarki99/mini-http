package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go listentoPprof()
	router := NewRouter()
	router.HandlerFunc("/test", func(w ResponseWriter, r *Request) {

		select {
		case <-r.Context.Done():
			fmt.Println("Request canceled by clie§nt during processing.")
			return
		default:
			w.WriteStatusCode(200)
			w.Header("Content-Type", "application/json")
			response := map[string]interface{}{
				"message": "Hello, World!",
			}
			data, _ := json.Marshal(response)
			fmt.Println("Response: ", string(data))
			w.Write(data)
		}

	})
	server := NewServer()
	fmt.Println("Starting to listen at 8080...")
	if err := server.StartListening("localhost:8080", router); err != nil {
		panic(err)
	}
}

func listentoPprof() {
	// Start a separate HTTP server for pprof
	pprofServer := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("pprof server running on :8081")
	if err := pprofServer.ListenAndServe(); err != nil {
		fmt.Printf("pprof server error: %v\n", err)
	}
}

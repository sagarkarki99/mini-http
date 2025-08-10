package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	router := NewRouter()
	router.HandlerFunc("/test", func(w ResponseWriter, r *Request) {
		m := map[string]interface{}{}
		json.NewDecoder(r.body).Decode(&m)
		fmt.Println("Body: ", m)
		for i := 0; i < 1; i++ {
			select {
			case <-r.Context.Done():
				fmt.Println("Request canceled by clie§nt during processing.")
				return
			case <-time.After(time.Second * 1):
				fmt.Println("Slept for", i+1, "second(s)")
			}
		}
		w.WriteStatusCode(200)
		w.Header("Content-Type", "application/json")
		response := map[string]interface{}{
			"message": "Hello, World!",
		}
		data, _ := json.Marshal(response)
		w.Write(data)

	})
	server := NewServer()
	fmt.Println("Starting to listen at 8080...")
	if err := server.StartListening("localhost:8080", router); err != nil {
		panic(err)
	}
}

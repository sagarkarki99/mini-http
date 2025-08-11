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

		select {
		case <-r.Context.Done():
			fmt.Println("Request canceled by clie§nt during processing.")
			return
		case <-time.After(time.Second * 1):
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

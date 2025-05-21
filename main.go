package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

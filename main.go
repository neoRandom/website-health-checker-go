package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello, world!")		
	})

	log.Println("Server starting at http://localhost:8080...")
	http.ListenAndServe(":8080", nil);
}

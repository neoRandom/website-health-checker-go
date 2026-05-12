package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, internet!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)

	fmt.Println("Starting server at http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}

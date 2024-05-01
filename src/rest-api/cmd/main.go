package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequest)

	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!", r.Body)
}

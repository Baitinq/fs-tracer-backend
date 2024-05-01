package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequest)

	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "Hello, World!", string(bytes))
}

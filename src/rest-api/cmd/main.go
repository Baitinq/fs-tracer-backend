package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello folks!")
	})
	mux.HandleFunc("/payload", handleRequest)

	http.ListenAndServe(":8080", mux)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "Hello, World!", string(bytes))
	log.Println("Request received", r.RemoteAddr, string(bytes))
}

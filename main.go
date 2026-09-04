package main

import "net/http"

const timeout = uint32(1)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /snapshot", webcamFrameHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

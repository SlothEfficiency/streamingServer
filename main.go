package main

import "net/http"

const timeout = uint32(1)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", serveFileHandler)

	mux.HandleFunc("GET /snapshot", webcamFrameHandler)
	mux.HandleFunc("GET /stream", webcamStreamHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/")
}

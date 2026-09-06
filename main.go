package main

import "net/http"

const timeout = uint32(1)

func main() {
	cam := NewCamera()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", serveFileHandler)

	mux.HandleFunc("GET /stream", cam.webcamStreamHandler)
	mux.HandleFunc("GET /state", cam.stateReaderHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/")
}

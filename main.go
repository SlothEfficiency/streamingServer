package main

import "net/http"

const timeout = uint32(1)

func main() {
	channels := NewChannelCollection()
	go channels.webcamMaster()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", serveFileHandler)

	mux.HandleFunc("GET /stream", channels.webcamStreamHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/")
}

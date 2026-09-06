package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, internalMsg string, statusCode int, err error) {
	if internalMsg != "" {
		log.Printf("%s: Sending Statuscode %v\n", internalMsg, statusCode)
	}

	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	byteResponse, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println("Error couldnt be processed.")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(byteResponse)
}

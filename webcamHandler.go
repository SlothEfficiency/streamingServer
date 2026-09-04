package main

import (
	"log"
	"net/http"

	"github.com/blackjack/webcam"
)

func webcamStreamHandler(w http.ResponseWriter, r *http.Request) {
	cam, err := initializeWebcam("Motion-JPEG")
	if err != nil {
		sendError(w, "Failed to initialize cam", 500, err)
	}
	defer cam.Close()

	sendResponse(w, 200, "multipart/x-mixed-replace; boundary=frame", []byte(""))

	for {
		select {
		case <-r.Context().Done():
			return
		default:
			frame, err := nextFrame(cam, timeout)
			if err != nil {
				sendError(w, "Failed to read frame", 500, err)
			}
			w.Write([]byte("--frame\r\nContent-Type: image/jpeg\r\n\r\n"))
			w.Write(frame)
			w.Write([]byte("\r\n"))
		}
	}
}

func webcamFrameHandler(w http.ResponseWriter, r *http.Request) {
	cam, err := initializeWebcam("Motion-JPEG")
	if err != nil {
		sendError(w, "Failed to initialize cam", 500, err)
	}
	defer cam.Close()

	frame, err := nextFrame(cam, timeout)
	if err != nil {
		sendError(w, "Failed to read frame", 500, err)
	}

	sendResponse(w, 200, "image/jpeg", frame)
}

func initializeWebcam(frameFormat string) (*webcam.Webcam, error) {
	// Initialize Camera
	cam, err := webcam.Open("/dev/video0") // Open webcam
	if err != nil {
		log.Println(err)
		return &webcam.Webcam{}, err
	}

	// Set Format
	err = setCamFormat(cam, frameFormat)
	if err != nil {
		log.Println(err)
		return &webcam.Webcam{}, err
	}
	// Start webcam
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return &webcam.Webcam{}, err
	}

	return cam, nil
}

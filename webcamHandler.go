package main

import (
	"log"
	"net/http"

	"github.com/blackjack/webcam"
)

type ChannelCollection struct {
	CamReader       chan []byte
	NewRequest      chan struct{}
	CloseConnection chan struct{}
	ErrorOccured    chan error
}

func NewChannelCollection() *ChannelCollection {
	return &ChannelCollection{
		CamReader:       make(chan []byte, 100),
		NewRequest:      make(chan struct{}),
		CloseConnection: make(chan struct{}),
		ErrorOccured:    make(chan error),
	}
}

func (col *ChannelCollection) webcamStreamHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	headerAlreadySet := false

	col.NewRequest <- struct{}{}

	sendError(w, "Failed to initialize cam", 500, err)

	for {
		select {

		// Tell webcamMaster that the connection is closed
		case <-r.Context().Done():
			col.CloseConnection <- struct{}{}

		// In case something goes wrong
		case err = <-col.ErrorOccured:
			if headerAlreadySet == false {
				sendError(w, "Something went wrong", 500, err)
				return
			}
			log.Println("Something went wrong, but it should not affect our connection.")
			continue

		// Read frame and sent it
		case frame := <-col.CamReader:
			if headerAlreadySet == false {
				w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
				w.WriteHeader(200)
				headerAlreadySet = true
			}
			w.Write([]byte("\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"))
			w.Write(frame)
			w.Write([]byte(""))
		}
	}
}

func initializeWebcam(frameFormat string) (*webcam.Webcam, error) {
	// Open Webcam
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		log.Println(err)
		return &webcam.Webcam{}, err
	}

	// Set Format
	err = setCamFormat(cam, frameFormat)
	if err != nil {
		cam.Close()
		log.Println(err)
		return &webcam.Webcam{}, err
	}
	return cam, nil
}

func (col *ChannelCollection) webcamMaster() {
	var err error

	cameraOpened := false
	OpenStreamsCounter := 0
	cam := &webcam.Webcam{}

	for {
		select {

		// New incoming request
		case <-col.NewRequest:
			// Start camera if it is the first one
			if OpenStreamsCounter == 0 {
				cam, err = initializeWebcam("Motion-JPEG")
				if err != nil {
					log.Println("Couldn't start camera because of ", err)
					col.ErrorOccured <- err
					continue
				}
				err = cam.StartStreaming()
				if err != nil {
					log.Println("Couldn't start stream because of ", err)
					cam.Close()
					col.ErrorOccured <- err
					continue
				}
				cameraOpened = true
			}
			OpenStreamsCounter += 1

		// Closed connection
		case <-col.CloseConnection:
			// The last one closes the door
			if OpenStreamsCounter == 1 {
				err = cam.Close()
				if err != nil {
					log.Println("Couldn't close camera because of ", err)
					continue
				}
				cameraOpened = false
			}
			OpenStreamsCounter -= 1

		// Generate new frame
		default:
			if cameraOpened {
				frame, err := nextFrame(cam, timeout)
				if err != nil {
					log.Printf("Couldn't read frame: %v", err)
					continue
				}
				col.CamReader <- frame
			}
		}
	}
}

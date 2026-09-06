package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blackjack/webcam"
)

type Camera struct {
	Cam                *webcam.Webcam
	CamReader          chan []byte
	OpenStreamsCounter int
	StopStream         chan struct{}
	mu                 *sync.Mutex
	StreamStopped      chan struct{}
}

func NewCamera() *Camera {
	return &Camera{
		CamReader:          make(chan []byte, 100),
		StopStream:         make(chan struct{}),
		mu:                 &sync.Mutex{},
		OpenStreamsCounter: 0,
		StreamStopped:      make(chan struct{}),
	}
}

func (cam *Camera) stateReaderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("%v", cam.OpenStreamsCounter)))
}

func (cam *Camera) webcamStreamHandler(w http.ResponseWriter, r *http.Request) {
	cam.mu.Lock()
	if cam.OpenStreamsCounter == 0 {
		cam.OpenStreamsCounter += 1
		err := cam.initializeWebcam("Motion-JPEG")
		if err != nil {
			sendError(w, "Failed to initialize cam", 500, err)
			return
		}
		cam.mu.Unlock()
		go cam.startStreaming()
		log.Println("Stream was started.")
	} else {
		cam.OpenStreamsCounter += 1
		cam.mu.Unlock()
	}

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.WriteHeader(200)

	for {
		select {
		case <-r.Context().Done():
			cam.mu.Lock()
			cam.OpenStreamsCounter -= 1
			if cam.OpenStreamsCounter == 0 {
				cam.StopStream <- struct{}{}
				<-cam.StreamStopped
			}
			cam.mu.Unlock()
			return
		case frame := <-cam.CamReader:
			w.Write([]byte("\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"))
			w.Write(frame)
			w.Write([]byte(""))
		}
	}
}

func (cam *Camera) initializeWebcam(frameFormat string) error {
	var err error

	cam.Cam, err = webcam.Open("/dev/video0")
	if err != nil {
		log.Println(err)
		return err
	}

	// Set Format
	err = setCamFormat(cam.Cam, frameFormat)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (cam *Camera) startStreaming() error {
	cam.mu.Lock()
	err := cam.Cam.StartStreaming()

	if err != nil {
		log.Println(err)
		return err
	}
	cam.mu.Unlock()

	for {
		cam.mu.Lock()
		select {
		case <-cam.StopStream:
			cam.Cam.Close()
			cam.StreamStopped <- struct{}{}
			return nil
		default:
			frame, err := nextFrame(cam.Cam, timeout)
			if err != nil {
				log.Printf("Couldn't read frame: %v", err)
			}
			cam.CamReader <- frame
		}
		cam.mu.Unlock()
	}
}

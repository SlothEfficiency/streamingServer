package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/blackjack/webcam"
)

type Camera struct {
	Cam                *webcam.Webcam
	CamReader          chan []byte
	OpenStreamsCounter atomic.Int32
	StopStream         chan struct{}
}

func NewCamera() *Camera {
	return &Camera{
		CamReader:  make(chan []byte, 100),
		StopStream: make(chan struct{}),
	}
}

func (cam *Camera) webcamStreamHandler(w http.ResponseWriter, r *http.Request) {
	if cam.OpenStreamsCounter.CompareAndSwap(0, 1) {
		err := cam.initializeWebcam("Motion-JPEG")
		if err != nil {
			sendError(w, "Failed to initialize cam", 500, err)
			cam.Cam.Close()
			return
		}
		go cam.startStreaming()
	}

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.WriteHeader(200)

	for {
		select {
		case <-r.Context().Done():
			cam.OpenStreamsCounter.Add(-1)
			if cam.OpenStreamsCounter.Load() == 0 {
				cam.StopStream <- struct{}{}
			}
			return
		case frame := <-cam.CamReader:
			w.Write([]byte("\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"))
			w.Write(frame)
			w.Write([]byte(""))
		}
	}
}

func (cam *Camera) initializeWebcam(frameFormat string) error {
	// Initialize Camera
	var err error
	cam.Cam, err = webcam.Open("/dev/video0") // Open webcam
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
	err := cam.Cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return err
	}
	for {
		select {
		case <-cam.StopStream:
			cam.Cam.Close()
			close(cam.CamReader)
			return nil
		default:
			frame, err := nextFrame(cam.Cam, timeout)
			if err != nil {
				log.Printf("Couldn't read frame: %v", err)
			}
			cam.CamReader <- frame
			log.Println("Neuer Frame gelesen.")
		}
	}
}

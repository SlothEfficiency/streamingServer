package main

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/blackjack/webcam"
)

type Camera struct {
	Cam                *webcam.Webcam
	CamReader          chan []byte
	OpenStreamsCounter atomic.Int32
	StopStream         chan struct{}
	mu                 *sync.Mutex
}

func NewCamera() *Camera {
	return &Camera{
		CamReader:  make(chan []byte, 100),
		StopStream: make(chan struct{}),
		mu:         &sync.Mutex{},
	}
}

func (cam *Camera) stateReaderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(string(cam.OpenStreamsCounter.Load())))
}

func (cam *Camera) webcamStreamHandler(w http.ResponseWriter, r *http.Request) {
	if cam.OpenStreamsCounter.CompareAndSwap(0, 1) {
		err := cam.initializeWebcam("Motion-JPEG")
		if err != nil {
			sendError(w, "Failed to initialize cam", 500, err)
			return
		}
		go cam.startStreaming()
	}

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.WriteHeader(200)

	for {
		select {
		case <-r.Context().Done():
			if cam.OpenStreamsCounter.CompareAndSwap(1, 0) {
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
	var err error
	cam.mu.Lock()
	defer cam.mu.Unlock()

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
	err := cam.Cam.StartStreaming()
	cam.mu.Lock()
	defer cam.mu.Unlock()
	if err != nil {
		log.Println(err)
		return err
	}

	for {
		select {
		case <-cam.StopStream:
			cam.Cam.Close()
			return nil
		default:
			frame, err := nextFrame(cam.Cam, timeout)
			if err != nil {
				log.Printf("Couldn't read frame: %v", err)
			}
			cam.CamReader <- frame
		}
	}
}

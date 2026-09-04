package main

import (
	"github.com/blackjack/webcam"
)

func nextFrame(cam *webcam.Webcam, timeout uint32) ([]byte, error) {
	err := cam.WaitForFrame(timeout)
	if err != nil {
		return nil, err
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		return nil, err
	}
	return frame, nil
}

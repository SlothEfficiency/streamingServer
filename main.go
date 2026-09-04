package main

import (
	"log"
	"os"

	"github.com/blackjack/webcam"
)

const timeout = uint32(1)

func main() {

	// Initialize Camera
	cam, err := webcam.Open("/dev/video0") // Open webcam
	if err != nil {
		log.Println(err)
		return
	}
	defer cam.Close()

	// Set Format
	err = setCamFormat(cam, "Motion-JPEG")
	if err != nil {
		log.Println(err)
		return
	}
	// Start webcam
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// Take a picture
		frame, err := nextFrame(cam, timeout)
		if err != nil {
			log.Println(err)
			return
		}

		//Save the picture
		err = os.WriteFile("frame.jpeg", frame, 0666)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

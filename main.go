package main

import (
	"fmt"
	"os"

	"github.com/blackjack/webcam"
)

func main() {

	// Initialize Camera
	cam, err := webcam.Open("/dev/video0") // Open webcam
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	// Set Format
	fmt.Println(w.Get)
	pixelformat, width, height, err := cam.SetImageFormat(1196444237, 1280, 720)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Format: %v, width: %v, height: %v\n", pixelformat, width, height)

	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}
	err = cam.WaitForFrame(20)

	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		fmt.Println("Timeout")
		fmt.Fprint(os.Stderr, err.Error())
	default:
		panic(err.Error())
	}

	frame, err := cam.ReadFrame()
	if len(frame) != 0 {
		os.WriteFile("frame.jpeg", frame, 0666)
	} else if err != nil {
		panic(err.Error())
	}
}

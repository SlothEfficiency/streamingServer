package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/blackjack/webcam"
)

func setCamFormat(cam *webcam.Webcam, formatName string) error {
	formats := cam.GetSupportedFormats()

	// Since the formats are in the keys and the formatNames are in the values, we have to compare the value with the name
	for format, name := range formats {
		if name == formatName {
			frameSizes := cam.GetSupportedFrameSizes(format)

			// Always choose max Resolution
			maxFrameSize := findNthBiggestFrameSize(frameSizes, 2)
			pixelformat, width, height, err := cam.SetImageFormat(format, maxFrameSize.MaxWidth, maxFrameSize.MaxHeight)
			if err != nil {
				fmt.Println("Couldn't set CamFormat.")
				return err
			}
			fmt.Printf("Format: %v, width: %v, height: %v\n", pixelformat, width, height)
			return nil
		}
	}

	// When we have no match return error
	return fmt.Errorf("CamFormat was not found.")
}

func findNthBiggestFrameSize(frameSizes []webcam.FrameSize, n uint32) webcam.FrameSize {
	log.Printf("Possible Framesizes: %v\n", frameSizes)
	if n > uint32(len(frameSizes)) {
		return frameSizes[len(frameSizes)-1]
	}

	slices.SortFunc(frameSizes, func(a, b webcam.FrameSize) int {
		return int(b.MaxWidth*b.MaxHeight - a.MaxWidth*a.MaxHeight)
	})
	return frameSizes[n-1]
}

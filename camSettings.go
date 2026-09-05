package main

import (
	"fmt"

	"github.com/blackjack/webcam"
)

func setCamFormat(cam *webcam.Webcam, formatName string) error {
	formats := cam.GetSupportedFormats()

	// Since the formats are in the keys and the formatNames are in the values, we have to compare the value with the name
	for format, name := range formats {
		if name == formatName {
			frameSizes := cam.GetSupportedFrameSizes(format)

			// Always choose max Resolution
			maxFrameSize := findMaxFrameSize(frameSizes)
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

func findMaxFrameSize(frameSizes []webcam.FrameSize) webcam.FrameSize {
	maxIndex := 0
	maxPixel := uint32(10000000)

	// iterate over formats to find highest resolution
	for i, frameSize := range frameSizes {
		if frameSize.MaxWidth*frameSize.MaxHeight < maxPixel {
			maxIndex = i
			maxPixel = frameSize.MaxWidth * frameSize.MaxHeight
		}
	}
	return frameSizes[maxIndex]
}

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
	length := len(frameSizes)
	if n == 0 || n > uint32(length) {
		// Return the smallest frame size or handle error as appropriate
		return frameSizes[length-1]
	}

	// Make a copy to avoid mutating the original slice
	sliceToSort := make([]webcam.FrameSize, length)
	copy(sliceToSort, frameSizes)

	log.Printf("Possible Framesizes: %v\n", sliceToSort)

	// Sort descending by area (MaxWidth * MaxHeight)
	slices.SortFunc(sliceToSort, func(a, b webcam.FrameSize) int {
		areaA := a.MaxWidth * a.MaxHeight
		areaB := b.MaxWidth * b.MaxHeight
		if areaA < areaB {
			return 1
		} else if areaA > areaB {
			return -1
		}
		return 0
	})
	log.Printf("Choosen %vth frameSize: %v\n", n, sliceToSort[n-1])
	log.Printf("Possible Framesizes: %v\n", sliceToSort)

	return sliceToSort[n-1]
}

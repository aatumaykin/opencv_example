// How to run:
//
// 		go run main.go 0
//
// +build example

package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

const MinimumArea = 3000

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tgo run main.go [cam ID]")
		return
	}

	deviceID := os.Args[1]

	cam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer cam.Close()

	window2 := gocv.NewWindow("imgDelta")
	window2.MoveWindow(800, 0)
	defer window2.Close()

	window3 := gocv.NewWindow("imgThresh")
	window3.MoveWindow(0, 900)
	defer window3.Close()

	window4 := gocv.NewWindow("imgDilate")
	window4.MoveWindow(800, 900)
	defer window4.Close()

	window := gocv.NewWindow("Motion Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	imgDilate := gocv.NewMat()
	defer imgDilate.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	status := "Ready"

	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := cam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		status = "Ready"
		statusColor := color.RGBA{0, 255, 0, 0}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(img, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgDilate, kernel)

		// now find contours
		contours := gocv.FindContours(imgDilate, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i, c := range contours {
			area := gocv.ContourArea(c)
			if area < MinimumArea {
				continue
			}

			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		window.IMShow(img)
		window2.IMShow(imgDelta)
		window3.IMShow(imgThresh)
		window4.IMShow(imgDilate)

		if gocv.WaitKey(1) >= 0 {
			break
		}
	}
}

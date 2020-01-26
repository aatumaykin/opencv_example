// How to run:
//
// 		go run main.go ../data/apple.jpg -threshold 128
//
// +build example

package main

import (
	"errors"
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"log"
)

var (
	thresh int
)

type ThreshMethods struct {
	name  string
	value gocv.ThresholdType
}

func createWindow(title string, position image.Point) *gocv.Window {
	win := gocv.NewWindow(title)
	win.MoveWindow(position.X, position.Y)

	return win
}

func readImage(filename string) gocv.Mat {
	log.Printf("Reading image from %v", filename)

	img := gocv.IMRead(filename, gocv.IMReadColor)

	if img.Empty() {
		log.Fatal(errors.New(fmt.Sprintf("Error reading image from: %v", filename)))
	}

	return img
}

func main() {
	flag.Usage = func() {
		log.Println("How to run:\n\tgo run main.go [image]")
		flag.PrintDefaults()
	}

	flag.IntVar(&thresh, "threshold", 128, "Threshold value")

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()

		return
	}

	inputs := flag.Args()

	file := inputs[0]

	img := readImage(file)
	defer img.Close()

	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)

	methods := []ThreshMethods{
		ThreshMethods{"THRESH_BINARY", gocv.ThresholdBinary},
		ThreshMethods{"THRESH_BINARY_INV", gocv.ThresholdBinaryInv},
		ThreshMethods{"THRESH_TRUNC", gocv.ThresholdTrunc},
		ThreshMethods{"THRESH_TOZERO", gocv.ThresholdToZero},
		ThreshMethods{"THRESH_TOZERO_INV", gocv.ThresholdToZeroInv},
	}

	winImg := createWindow("Image", image.Point{})
	winImg.IMShow(img)

	gocv.WaitKey(0)

	for _, method := range methods {
		threshImg := gocv.NewMat()
		gocv.Threshold(img, &threshImg, float32(thresh), 255, method.value)

		winImg.SetWindowTitle(method.name)
		winImg.IMShow(threshImg)

		gocv.WaitKey(0)
	}
}

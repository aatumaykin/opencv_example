// How to run:
//
// 		go run main.go ../data/shapes_example.png
//
// +build example

package main

import (
	"errors"
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
)

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

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()

		return
	}

	inputs := flag.Args()

	file := inputs[0]

	img := readImage(file)

	winImg := createWindow("Image", image.Point{})
	winImg.IMShow(img)

	mask := gocv.NewMat()
	gocv.InRangeWithScalar(img, gocv.NewScalar(0, 0, 0, 0), gocv.NewScalar(15, 15, 15, 0), &mask)

	winMask := createWindow("Mask", image.Point{})
	winMask.IMShow(mask)

	gocv.WaitKey(0)

	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	fmt.Printf("I found %d black shapes", len(contours))

	for i, _ := range contours {
		gocv.DrawContours(&img, contours, i, color.RGBA{150, 150, 255, 0}, 3)
	}

	winImg.IMShow(img)

	gocv.WaitKey(0)
}

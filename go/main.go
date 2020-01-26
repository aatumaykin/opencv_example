// How to run:
//
// 		go run main.go [image]
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

	// define the list of boundaries
	boundaries := [][]gocv.Scalar{
		{gocv.NewScalar(15, 40, 140, 0), gocv.NewScalar(50, 90, 200, 0)},
		{gocv.NewScalar(130, 60, 15, 0), gocv.NewScalar(180, 80, 50, 0)},
		{gocv.NewScalar(0, 118, 130, 0), gocv.NewScalar(62, 150, 170, 0)},
		{gocv.NewScalar(40, 40, 20, 0), gocv.NewScalar(70, 50, 50, 0)},
	}

	img := readImage(file)
	defer img.Close()

	winImg := createWindow("Image", image.Point{})
	winImg.IMShow(img)

	gocv.WaitKey(0)

	winImg.Close()

	winOut := createWindow("Image", image.Point{})

	for _, bound := range boundaries {
		mask := gocv.NewMat()
		gocv.InRangeWithScalar(img, bound[0], bound[1], &mask)
		gocv.Merge([]gocv.Mat{mask,mask,mask}, &mask)

		output := gocv.NewMat()
		gocv.BitwiseAnd(img, mask, &output)

		winOut.IMShow(output)
		gocv.WaitKey(0)
	}
}

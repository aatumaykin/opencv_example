// How to run:
//
// 		go run main.go ../data/shapes_and_colors.jpg
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

	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)

	imgBlurred := gocv.NewMat()
	gocv.GaussianBlur(imgGray, &imgBlurred, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	imgThresh := gocv.NewMat()
	gocv.Threshold(imgBlurred, &imgThresh, 60, 255, gocv.ThresholdBinary)

	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	for i, contour := range contours {
		gocv.DrawContours(&img, contours, i, color.RGBA{0, 255, 0, 0}, 2)

		mat := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
		gocv.FillPoly(&mat, [][]image.Point{contour}, color.RGBA{255, 255, 255, 1})

		moments := gocv.Moments(mat, true)

		x := int(moments["m10"] / moments["m00"])
		y := int(moments["m01"] / moments["m00"])

		gocv.Circle(&img, image.Pt(x, y), 7, color.RGBA{255, 255, 255, 0}, -1)
	}

	winImg := createWindow("Image", image.Point{})
	winImg.IMShow(img)

	gocv.WaitKey(0)
}

// How to run:
//
// 		go run main.go ../data/tetris_blocks.png
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

var (
	contourColor = color.RGBA{240, 0, 159, 0}

	// windows
	winImg        *gocv.Window
	winGray       *gocv.Window
	winEdge       *gocv.Window
	winThresh     *gocv.Window
	winContours   *gocv.Window
	winErode      *gocv.Window
	winDilate     *gocv.Window
	winBitwiseAnd *gocv.Window
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
	defer img.Close()

	winImg = createWindow("Image", image.Point{})

	winGray = createWindow("Gray", image.Pt(500, 0))
	defer winGray.Close()

	winEdge = createWindow("Edge", image.Pt(1000, 0))
	defer winEdge.Close()

	winThresh = createWindow("Thresh", image.Pt(0, 600))
	defer winThresh.Close()

	winContours = createWindow("Contours", image.Pt(500, 600))
	defer winContours.Close()

	winErode = createWindow("Erode", image.Pt(1000, 600))
	defer winErode.Close()

	winDilate = createWindow("Dilate", image.Pt(0, 900))
	defer winDilate.Close()

	winBitwiseAnd = createWindow("BitwiseAnd", image.Pt(500, 900))
	defer winBitwiseAnd.Close()

	// convert the image to grayscale
	imgGray := gocv.NewMat()
	defer imgGray.Close()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	defer kernel.Close()

	// applying edge detection we can find the outlines of objects in images
	imgEdges := gocv.NewMat()
	defer imgEdges.Close()
	gocv.Canny(imgGray, &imgEdges, 30, 150)

	// threshold the image by setting all pixel values less than 225
	// to 255 (white; foreground) and all pixel values >= 225 to 255
	// (black; background), thereby segmenting the image
	imgThresh := gocv.NewMat()
	defer imgThresh.Close()
	gocv.Threshold(imgGray, &imgThresh, 225, 255, gocv.ThresholdBinaryInv)

	imgContours := img.Clone()

	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	for i, _ := range contours {
		gocv.DrawContours(&imgContours, contours, i, contourColor, 3)
	}

	// draw the total number of contours found in purple
	gocv.PutText(&imgContours, fmt.Sprintf("Found: %d objects", len(contours)), image.Pt(10, 20), gocv.FontHersheySimplex, 0.7, color.RGBA{240, 0, 159, 0}, 2)

	// we apply erosions to reduce the size of foreground objects
	imgErode := gocv.NewMat()
	defer imgErode.Close()
	gocv.Erode(imgThresh, &imgErode, kernel)

	// similarly, dilations can increase the size of the ground objects
	imgDilate := gocv.NewMat()
	defer imgDilate.Close()
	gocv.Dilate(imgThresh, &imgDilate, kernel)

	// a typical operation we may want to apply is to take our mask and
	// apply a bitwise AND to our input image, keeping only the masked
	// regions
	imgBitwiseAnd := gocv.NewMat()
	defer imgBitwiseAnd.Close()
	gocv.Merge([]gocv.Mat{imgThresh, imgThresh, imgThresh}, &imgBitwiseAnd)
	gocv.BitwiseAnd(img, imgBitwiseAnd, &imgBitwiseAnd)

	winImg.IMShow(img)
	winGray.IMShow(imgGray)
	winEdge.IMShow(imgEdges)
	winThresh.IMShow(imgThresh)
	winContours.IMShow(imgContours)
	winErode.IMShow(imgErode)
	winDilate.IMShow(imgDilate)
	winBitwiseAnd.IMShow(imgBitwiseAnd)

	for {
		if gocv.WaitKey(1) >= 0 {
			break
		}
	}
}

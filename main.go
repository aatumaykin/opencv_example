package main

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
)

const MinimumArea = 100

func viewImage(image gocv.Mat, windowName string) *gocv.Window {
	window := gocv.NewWindow(windowName)
	window.ResizeWindow(640, 480)
	window.IMShow(image)

	return window
}

// рисуем контуры
func drawContours(contours [][]image.Point, dest gocv.Mat, color color.RGBA) {
	for _, contour := range contours {
		if area := gocv.ContourArea(contour); area > MinimumArea {
			rect := gocv.BoundingRect(contour)
			gocv.Rectangle(&dest, rect, color, 2)
		}
	}
}

func printPixelColor(img gocv.Mat, point image.Point) {
	// split image channels
	bgr := gocv.Split(img)
	// pixel values for each channel - we know this is a BGR image
	fmt.Printf("Pixel B: %d\n", bgr[0].GetUCharAt(point.Y, point.X))
	fmt.Printf("Pixel G: %d\n", bgr[1].GetUCharAt(point.Y, point.X))
	fmt.Printf("Pixel R: %d\n", bgr[2].GetUCharAt(point.Y, point.X))
}

func readImage(filename string) gocv.Mat {
	log.Printf("Reading image from %v", filename)

	img := gocv.IMRead(filename, gocv.IMReadColor)

	if img.Empty() {
		log.Fatal(errors.New(fmt.Sprintf("Error reading image from: %v", filename)))
	}

	// size of an image
	fmt.Printf("%s size: %d x %d\n", filename, img.Rows(), img.Cols())
	// image channels
	fmt.Printf("%s channels: %d\n", filename, img.Channels())

	return img
}

func main() {
	img1 := readImage("apple.jpg")
	defer img1.Close()

	//img1 = crop(img1, image.Pt(0, 50), image.Pt(width, height))

	// Convert BGR to HSV image
	hsvImg1 := img1.Clone()
	gocv.CvtColor(img1, &hsvImg1, gocv.ColorBGRToHSV)

	// Create the mask
	lowerBound := gocv.NewMatFromScalar(gocv.NewScalar(0.0, 0.0, 255.0, 0.0), gocv.MatTypeCV8U)
	upperBound := gocv.NewMatFromScalar(gocv.NewScalar(255.0, 255.0, 255.0, 0.0), gocv.MatTypeCV8U)

	mask1 := gocv.NewMat()
	gocv.InRange(hsvImg1, lowerBound, upperBound, &mask1)

	// maskedImg: output array that has the same size and type as the input arrays.
	maskedImg1 := gocv.NewMatWithSize(hsvImg1.Rows(), hsvImg1.Cols(), gocv.MatTypeCV8U)
	hsvImg1.CopyToWithMask(&maskedImg1, mask1)

	// Create the inverted mask
	maskInv1 := gocv.NewMat()
	gocv.BitwiseNot(mask1, &maskInv1)

	// Convert to grayscale image
	gray1 := gocv.NewMat()
	gocv.CvtColor(img1, &gray1, gocv.ColorBGRToGray)

	// Bitwise-OR mask and original image
	gocv.BitwiseOr(img1, maskedImg1, &mask1)

	img1 = img1.Crop(image.Rect(0, 0, 50, 50))

	viewImage(img1, "Image 1")
	viewImage(hsvImg1, "hsvImg1").MoveWindow(50, 50)
	viewImage(mask1, "mask1").MoveWindow(100, 100)
	viewImage(maskedImg1, "maskedImg1").MoveWindow(150, 150)
	viewImage(maskInv1, "maskInv1").MoveWindow(200, 200)
	viewImage(gray1, "gray1").MoveWindow(250, 250)

	//viewImage(img2, "Image 2").MoveWindow(100, 100)

	printPixelColor(img1, image.Pt(100, 100))

	for {
		ch := gocv.WaitKey(1)

		if ch == 27 || ch == int('q') {
			break
		}
	}
}

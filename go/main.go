// How to run:
//
// 		go run main.go
//
// +build example

package main

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"strconv"
)

const MinimumArea = 300

var debug bool = false

type RegionData struct {
	region image.Rectangle
	tresh  float32
}

func viewImage(image gocv.Mat, windowName string) *gocv.Window {
	window := gocv.NewWindow(windowName)
	window.SetWindowProperty(gocv.WindowPropertyVisible, gocv.WindowNormal)
	if !image.Empty() {
		window.IMShow(image)
	}

	return window
}

// рисуем контуры
func drawContours(contours [][]image.Point, dest gocv.Mat, color color.RGBA) {
	for _, contour := range contours {
		if area := gocv.ContourArea(contour); area > MinimumArea {
			rect := gocv.BoundingRect(contour)
			gocv.Rectangle(&dest, rect, color, 3)
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

func getRegionsData() []RegionData {
	var data []RegionData
	maxWidth := 1920
	maxHeight := 1080

	// зона 1
	data = append(data, RegionData{
		image.Rectangle{Min: image.Pt(1560, 920), Max: image.Pt(maxWidth, maxHeight)},
		18.0,
	})

	// зона 2
	data = append(data, RegionData{
		image.Rectangle{Min: image.Pt(610, 900), Max: image.Pt(760, maxHeight)},
		12.0,
	})

	// зона 3
	data = append(data, RegionData{
		image.Rectangle{Min: image.Pt(1050, 590), Max: image.Pt(1170, 660)},
		22.0,
	})

	// зона 4
	data = append(data, RegionData{
		image.Rectangle{Min: image.Pt(1780, 620), Max: image.Pt(maxWidth, 760)},
		12.0,
	})

	// зона 5
	data = append(data, RegionData{
		image.Rectangle{Min: image.Pt(1600, 540), Max: image.Pt(1779, 740)},
		20.0,
	})

	return data
}

func drawRegion(index int, regionsData []RegionData, imgOrigin, imgDetected gocv.Mat) {
	regionOrigin := imgOrigin.Region(regionsData[index].region)
	regionDetected := imgDetected.Region(regionsData[index].region)

	mask := gocv.NewMat()
	gocv.AbsDiff(regionOrigin, regionDetected, &mask)
	gocv.CvtColor(mask, &mask, gocv.ColorBGRToGray)

	// apply Gaussian blur
	gocv.GaussianBlur(mask, &mask, image.Pt(11, 11), 0, 0, gocv.BorderDefault)

	// create binary image
	threshImg := gocv.NewMat()
	gocv.Threshold(mask, &threshImg, regionsData[index].tresh, 255.0, gocv.ThresholdBinary)

	contours := gocv.FindContours(threshImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	gocv.Rectangle(&imgDetected, regionsData[index].region, color.RGBA{G: 255}, 1)
	drawContours(contours, regionDetected, color.RGBA{R: 255})


	if debug {
		drawContours(contours, mask, color.RGBA{R: 255, G: 255, B: 255})
		drawContours(contours, threshImg, color.RGBA{R: 255, G: 255, B: 255})

		viewImage(regionDetected, "Region-"+strconv.Itoa(index)).MoveWindow(1000, 0)
		viewImage(mask.Clone(), "Mask-"+strconv.Itoa(index)).MoveWindow(1000, regionDetected.Rows()+50)
		viewImage(threshImg.Clone(), "Thresh-"+strconv.Itoa(index)).MoveWindow(1000, regionDetected.Rows()+mask.Rows()+50)
	}
}

func main() {
	imgOrigin := readImage("../data/origin-day.jpg")
	defer imgOrigin.Close()

	files := []string{
		"../data/A20012713182810.jpg", // cat 0
		"../data/A20012713203610.jpg", // cat 1
		"../data/A20012713203711.jpg", // cat 2
		"../data/A20012713203812.jpg", // cat 3
		"../data/A20012713193210.jpg", // cat 4
		"../data/A20012713193311.jpg", // cat 5
		"../data/A20012713193412.jpg", // cat 6
		"../data/A20012713480211.jpg", // cat 7
		"../data/A20012713480412.jpg", // cat 8
		"../data/A20012808501810.jpg", // cat 9
		"../data/A20012808501911.jpg", // cat 10
		"../data/A20012808502012.jpg", // cat 11
	}

	fileIndex := 10

Loop:
	for {
		if fileIndex >= len(files) {
			fileIndex = 0
		}

		if fileIndex < 0 {
			fileIndex = len(files) - 1
		}

		imgDetected := readImage(files[fileIndex])
		defer imgDetected.Close()

		printPixelColor(imgOrigin, image.Pt(100, 100))

		//copyImgDetected := imgDetected.Clone()

		regionsData := getRegionsData()
		for i, _ := range regionsData {
			drawRegion(i, regionsData, imgOrigin, imgDetected)
		}

		win := viewImage(imgDetected, "Image")
		win.SetWindowTitle(strconv.Itoa(fileIndex) + " - " + files[fileIndex])
		defer win.Close()

		thresh := float32(0.0)

		//ind := 3
		//regionOrigin := imgOrigin.Region(regionsData[ind].region)
		//regionDetected := copyImgDetected.Region(regionsData[ind].region)
		//
		//viewImage(regionOrigin.Clone(), "regionOrigin")
		//viewImage(regionDetected.Clone(), "regionDetected")
		//
		//mask := gocv.NewMat()
		//gocv.AbsDiff(regionOrigin, regionDetected, &mask)
		//
		//viewImage(mask.Clone(), "diff")
		//
		//gocv.CvtColor(mask, &mask, gocv.ColorBGRToGray)
		//viewImage(mask.Clone(), "color")
		//
		//// apply Gaussian blur
		//gocv.GaussianBlur(mask, &mask, image.Pt(11, 11), 0, 0, gocv.BorderDefault)
		//viewImage(mask.Clone(), "blur")
		//
		//winThresh := viewImage(gocv.NewMat(), "Thresh: ")
		//winThresh.MoveWindow(300, 600)

		for {
			//threshImg := gocv.NewMat()
			//gocv.Threshold(mask, &threshImg, thresh, 255.0, gocv.ThresholdBinary)
			//
			//contours := gocv.FindContours(threshImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)
			//
			//drawContours(contours, threshImg, color.RGBA{R: 255, G: 255, B: 255})
			//
			//winThresh.IMShow(threshImg)
			//winThresh.SetWindowTitle("Thresh: " + strconv.Itoa(int(thresh)))

			ch := gocv.WaitKey(10)

			if ch == 27 || ch == int('q') {
				break Loop
			}

			// right
			if ch == 3 {
				fileIndex++

				break
			}

			// left
			if ch == 2 {
				fileIndex--

				break
			}

			if ch == 32 {
				thresh += 1.0
			}
		}
	}
}

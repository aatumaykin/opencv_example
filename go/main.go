// How to run:
//
// 		go run main.go -gray -play -stopDetected -sleep 10 -start 0 ./images/
//
// +build example

package main

import (
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"time"
)

const MinimumArea = 1200
const MinimumPercent = 30

var (
	contourColor = color.RGBA{255, 0, 0, 0}
	rectColor    = color.RGBA{0, 0, 255, 0}
	statusColor  = color.RGBA{0, 255, 0, 0}

	cropRect = image.Rect(0, 50, 1920, 1080)

	currentPosition int
	gray         bool
	play         bool
	stopDetected bool
	sleep        int
)

func findContours(img gocv.Mat) [][]image.Point {
	finded := gocv.FindContours(img, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var contours = [][]image.Point{}

	for _, c := range finded {
		area := gocv.ContourArea(c)

		if area < MinimumArea {
			continue
		}

		rect := gocv.BoundingRect(c)

		regionImg := img.Region(rect)
		size := regionImg.Size()
		percent := math.Ceil(float64(gocv.CountNonZero(regionImg)) / float64(size[0]*size[1]) * 100)
		ratioY := float64(size[0]) / float64(size[1])
		ratioX := float64(size[1]) / float64(size[0])

		if percent < MinimumPercent || ratioY > 4 || ratioX > 4 {
			continue
		}

		fmt.Println(size, area, percent, fmt.Sprintf("%.2f %.2f", ratioX, ratioY))

		gocv.Rectangle(&img, rect, color.RGBA{255, 255, 255, 0}, 1)

		contours = append(contours, c)
	}

	return contours
}

func drawContours(contours [][]image.Point, img gocv.Mat) {
	for i, c := range contours {
		gocv.DrawContours(&img, contours, i, contourColor, 1)

		rect := gocv.BoundingRect(c)
		gocv.Rectangle(&img, rect, rectColor, 1)
	}
}

func offsetContours(contours [][]image.Point, offset image.Point) [][]image.Point {
	tmp := contours

	for i, c := range tmp {
		for j, p := range c {
			tmp[i][j] = p.Add(offset)
		}
	}

	return tmp
}

func findFiles(root string) []string {
	files := []string{}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".jpg" {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func createWindow(title string, position image.Point) *gocv.Window {
	win := gocv.NewWindow(title)
	win.MoveWindow(position.X, position.Y)

	return win
}

func excludeArea(img gocv.Mat) {
	triangle := [][]image.Point{{image.Pt(525, 0), image.Pt(1920, 0), image.Pt(1920, 520)}}

	gocv.FillPoly(&img, triangle, color.RGBA{255, 255, 255, 0})
}

func main() {
	flag.Usage = func() {
		fmt.Println("How to run:\n\tgo run main.go [flag] [folder]")
		flag.PrintDefaults()
	}

	flag.BoolVar(&gray, "gray", false, "Use gray color space")
	flag.BoolVar(&play, "play", false, "Auto play or manual")
	flag.BoolVar(&stopDetected, "stopDetected", false, "Stop if detected")
	flag.IntVar(&sleep, "sleep", 1, "Sleep in ms between images")
	flag.IntVar(&currentPosition, "start", 0, "Index start image")

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()

		return
	}

	inputs := flag.Args()

	filepath := inputs[0]

	files := findFiles(filepath)

	window := createWindow("Motion Window", image.Pt(0, 0))
	defer window.Close()

	status := "Ready"

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()
Loop:
	for {
		img := gocv.IMRead(files[currentPosition], gocv.IMReadColor)

		if img.Empty() {
			break
		}

		status = fmt.Sprintf("Frame: %d. Ready", currentPosition)

		imgDelta := gocv.NewMat()
		imgThresh := gocv.NewMat()
		imgDilate := gocv.NewMat()

		cropImg := img.Clone()
		excludeArea(cropImg)
		cropImg = cropImg.Crop(cropRect)

		imgGray := cropImg.Clone()

		cropImg.Close()

		if gray {
			gocv.CvtColor(imgGray, &imgGray, gocv.ColorBGRToGray)
		}

		imgBlur := gocv.NewMat()
		gocv.Blur(imgGray, &imgBlur, image.Pt(9, 9))

		imgGray.Close()

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(imgBlur, &imgDelta)

		imgBlur.Close()

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		imgDelta.Close()

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgDilate, kernel)

		imgThresh.Close()

		contours := findContours(imgDilate)

		imgDilate.Close()

		contours = offsetContours(contours, image.Pt(0, 50))
		drawContours(contours, img)

		if len(contours) > 0 {
			status = fmt.Sprintf("Frame: %d. Motion detected", currentPosition)

			if stopDetected {
				play = false
			}
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		window.IMShow(img)

		if !play {
			for {
				ch := gocv.WaitKey(1)

				// esc || q
				if ch == 27 || ch == int('q') {
					break Loop
				}

				// left
				if ch == 2 {
					currentPosition--
					break
				}

				// right
				if ch == 3 {
					currentPosition++
					break
				}

				// space
				if ch == 32 {
					currentPosition++
					play = true
					break
				}
			}
		} else {
			currentPosition++

			ch := gocv.WaitKey(1)

			// esc || q
			if ch == 27 || ch == int('q') {
				break
			}

			// p
			if ch == int('p') {
				play = false
			}

			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
}

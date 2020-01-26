// How to run:
//
// 		go run main.go [flag] [video]
//
// +build example

package main

import (
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"time"
)

const MinimumArea = 1000

var (
	contourColor = color.RGBA{255, 0, 0, 0}
	rectColor    = color.RGBA{0, 0, 255, 0}
	statusColor  = color.RGBA{0, 255, 0, 0}

	currentPosition int = 0

	// flags
	play       bool
	sleep      int
	showDelta  bool
	showDilate bool
	showTresh  bool
	showGray   bool
	gray       bool

	// windows
	winFrame  *gocv.Window
	winDelta  *gocv.Window
	winDilate *gocv.Window
	winTresh  *gocv.Window
	winGray   *gocv.Window
)

func createWindow(title string, position image.Point) *gocv.Window {
	win := gocv.NewWindow(title)
	win.MoveWindow(position.X, position.Y)

	return win
}

func findContours(img gocv.Mat) [][]image.Point {
	var contours = [][]image.Point{}

	allContours := gocv.FindContours(img, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	for _, contour := range allContours {
		area := gocv.ContourArea(contour)

		if area < MinimumArea {
			continue
		}

		contours = append(contours, contour)
	}

	return contours
}

func drawContours(contours [][]image.Point, img gocv.Mat) {
	for i, contour := range contours {
		gocv.DrawContours(&img, contours, i, contourColor, 1)

		rect := gocv.BoundingRect(contour)
		gocv.Rectangle(&img, rect, rectColor, 1)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("How to run:\n\tgo run main.go [flag] [video]")
		flag.PrintDefaults()
	}

	flag.BoolVar(&play, "play", false, "Auto play or manual")
	flag.IntVar(&sleep, "sleep", 1, "Sleep in ms between images")
	flag.BoolVar(&showDelta, "showDelta", false, "Show Delta image")
	flag.BoolVar(&showDilate, "showDilate", false, "Show Dilate image")
	flag.BoolVar(&showTresh, "showTresh", false, "Show Tresh image")
	flag.BoolVar(&showGray, "showGray", false, "Show gray image")
	flag.BoolVar(&gray, "gray", false, "Use gray color space")

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()

		return
	}

	inputs := flag.Args()

	file := inputs[0]

	video, err := gocv.VideoCaptureFile(file)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", file)
		return
	}
	defer video.Close()

	frames := int(video.Get(gocv.VideoCaptureFrameCount))
	fmt.Println("frames count %d", frames)

	if showDelta {
		winDelta = createWindow("Delta", image.Pt(800, 0))
		defer winDelta.Close()
	}

	if showTresh {
		winTresh = createWindow("Thresh", image.Pt(0, 900))
		defer winTresh.Close()
	}

	if showDilate {
		winDilate = createWindow("Dilate", image.Pt(800, 900))
		defer winDilate.Close()
	}

	if showGray {
		winGray = createWindow("Gray", image.Pt(900, 900))
		defer winGray.Close()
	}

	winFrame := createWindow("Motion Window", image.Point{})
	defer winFrame.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	imgDilate := gocv.NewMat()
	defer imgDilate.Close()

	imgGray := gocv.NewMat()
	defer imgGray.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	status := "Ready"

Loop:
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}

		if img.Empty() {
			continue
		}

		status = fmt.Sprintf("Frame: %d. Ready", currentPosition)

		imgGray := img.Clone()

		if gray {
			gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)
		}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(imgGray, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 125, 255, gocv.ThresholdBinary)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgDilate, kernel)

		// now find contours
		contours := findContours(imgDilate)
		drawContours(contours, img)

		if len(contours) > 0 {
			status = fmt.Sprintf("Frame: %d. Motion detected", currentPosition)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		winFrame.IMShow(img)

		if showDelta {
			winDelta.IMShow(imgDelta)
		}

		if showTresh {
			winTresh.IMShow(imgThresh)
		}

		if showDilate {
			winDilate.IMShow(imgDilate)
		}

		if showGray {
			winGray.IMShow(imgGray)
		}

		if !play {
			for {
				ch := gocv.WaitKey(1)

				// esc || q
				if ch == 27 || ch == int('q') {
					break Loop
				}

				// right
				if ch == 3 {
					currentPosition++;
					break
				}
			}
		} else {
			currentPosition++;
			if gocv.WaitKey(1) >= 0 {
				break
			}

			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
}

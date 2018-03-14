package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/trace"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
)

var (
	samples = 75
	width   = 600
	height  = 500
	numCPU  = runtime.NumCPU()
)

// Check is used for handling errors.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func saveFile(fileName string, img image.Image) error {
	fileName, err := filepath.Abs("../output/" + fileName)
	check(err)

	// If the file format is supported we create the file and
	// write the data to the file.
	if strings.Contains(fileName, ".png") {
		// Create the file if it doesn't exist already.
		file, err := os.Create(fileName)
		defer file.Close()
		check(err)

		return png.Encode(file, img)
	} else if strings.Contains(fileName, ".jpg") {
		// Create the file if it doesn't exist already.
		file, err := os.Create(fileName)
		defer file.Close()
		check(err)

		// Quality from 0-100.
		o := jpeg.Options{Quality: 100}
		return jpeg.Encode(file, img, &o)
	} else if strings.Contains(fileName, ".bmp") {
		// Create the file if it doesn't exist already.
		file, err := os.Create(fileName)
		defer file.Close()
		check(err)

		return bmp.Encode(file, img)
	}

	// No supported file format found.
	return fmt.Errorf("file format not supported, use: png, bmp or jpg")
}

func pixel(x, y int, scn *scene) color.NRGBA {
	// Starting point for each pixel.
	col := vec3{0.0, 0.0, 0.0}
	for i := 0; i < samples; i++ {
		// Add a bit of randomness, so the background will blend more with the edges of objects.
		// This will prevent lines from looking jaggy.
		u := (float64(x) + rand.Float64()) / float64(width)
		v := (float64(y) + rand.Float64()) / float64(height)

		r := scn.cam.getRay(u, v)
		col = col.add(r.color(scn, 0))
	}
	// Divide by the amount of samples to get the average.
	col = col.divScalar(float64(samples))
	return color.NRGBA{
		uint8(col.x * 255.0),
		uint8(col.y * 255.0),
		uint8(col.z * 255.0),
		255,
	}
}

func render(scn *scene) image.Image {
	fmt.Println("Number of samples:", samples)

	// The surface or rectangle for the image.
	rect := image.Rect(0, 0, width, height)
	rgba := image.NewNRGBA(rect)

	// Loop through each pixel from left to write. cx and cy being the current x and y respectively.
	for cy := 0; cy < height; cy++ {
		for cx := 0; cx < width; cx++ {
			rgba.SetNRGBA(cx, cy, pixel(cx, cy, scn))
		}
	}
	return rgba
}

func main() {
	// Creates a trace file with CPU usage and stuff.
	trcName, err := filepath.Abs("../trace.out")
	check(err)

	trcFile, err := os.Create(trcName)
	check(err)
	err = trace.Start(trcFile)
	check(err)
	defer trcFile.Close()
	defer trace.Stop()

	// Let Go use all the power the PC has.
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("Number of threads available:", numCPU)

	// Image dimensions
	fmt.Println("Image width:", width)
	fmt.Println("Image height:", height)

	// List of objects.
	objList := []object{
		object{shapeCircle, 0.5, v(-0.25, 0.0, -1.0), met(1.0, 0.0, 1.0, 0.0)},
		object{shapeCircle, 0.5, v(0.25, 0.0, -2.2), met(0.0, 1.0, 1.0, 0.2)},
		object{shapeCircle, 100.0, v(0.0, -100.5, -1.0), dif(0.5, 0.5, 0.5)},
	}

	// Create a scene, containing a camera and a list of objects to render.
	scn := scene{cam(v(0.0, 0.0, 1.0), v(0.0, 0.0, -1.0), 75.0, float64(width)/float64(height)), objList}

	// Get the current time, use this to get the elapsed time later.
	startTime := time.Now()

	img := render(&scn)

	// Print how long it took to raycast.
	elapsed := time.Since(startTime)
	fmt.Println("Time spent raycasting:", elapsed.Seconds(), "s")

	// Flip the image because we want 0, 0 to be the bottom left.
	img = imaging.FlipV(img)

	// Check if we have enough arguments, if not tell the user he should pass a file name.
	if len(os.Args) < 2 {
		err := fmt.Errorf("not enough arguments, usage:\n render test.png")
		check(err)
	}

	// Save the file to the destination given in the argument.
	err = saveFile(os.Args[1], img)
	check(err)
}

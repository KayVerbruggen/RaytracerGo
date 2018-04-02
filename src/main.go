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
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
)

var (
	samples = 100
	width   = 1000
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

func render(scn *scene) image.Image {
	fmt.Println("Number of samples:", samples)

	// The surface or rectangle for the image.
	rect := image.Rect(0, 0, width, height)
	rgba := image.NewNRGBA(rect)

	// Create goroutines for each row.
	var w sync.WaitGroup
	w.Add(height)

	// Loop through each pixel from left to write. cx and cy being the current x and y respectively.
	for cy := 0; cy < height; cy++ {
		go func(cy int) {
			// Create a rand.Rand interface for each goroutine to prevent locking and unlocking.
			// Give it a seed by generating a random number, time can be used but this gave a weird effect.
			rnd := rand.New(rand.NewSource(rand.Int63()))

			for cx := 0; cx < width; cx++ {
				// Starting point for each pixel.
				col := vec3{0.0, 0.0, 0.0}
				for i := 0; i < samples; i++ {
					// Add a bit of randomness, so the background will blend more with the edges of objects.
					// This will prevent lines from looking jaggy.
					s := (float64(cx) + rnd.Float64()) / float64(width)
					t := (float64(cy) + rnd.Float64()) / float64(height)

					r := scn.cam.ray(s, t, rnd)
					col = col.add(r.color(scn, 0, rnd))
				}
				// Divide by the amount of samples to get the average.
				col = col.divScalar(float64(samples))
				finalCol := color.NRGBA{
					uint8(col.x * 255.0),
					uint8(col.y * 255.0),
					uint8(col.z * 255.0),
					255,
				}
				rgba.SetNRGBA(cx, cy, finalCol)
			}
			w.Done()
		}(cy)
	}
	w.Wait()
	return rgba
}

func main() {
	// Check if we have enough arguments, if not tell the user he should pass a file name.
	if len(os.Args) < 2 {
		err := fmt.Errorf("not enough arguments, usage:\n render test.png")
		panic(err)
	}

	// Creates a trace file with CPU usage and stuff.
	trcName, err := filepath.Abs("../trace.out")
	check(err)
	trcFile, err := os.Create(trcName)
	check(err)

	// Start tracing.
	err = trace.Start(trcFile)
	check(err)

	// Close everything add the end of the program.
	defer trcFile.Close()
	defer trace.Stop()

	// Let the user know how many threads it is using.
	fmt.Println("Number of threads available:", numCPU)

	// Image dimensions
	fmt.Println("Image width:", width)
	fmt.Println("Image height:", height)
	/*
	   // List of objects.
	   objList := []*object{
	   	sphere(1000.0, v(0.0, -1000.0, 0.0), dif(0.5, 0.5, 0.5)),
	   	movingSphere(1.0, v(0.0, 1.0, 0.0), v(0.0, 1+0.5*rand.Float64(), 0.0), 0.0, 1.0, dif(1.0, 0.2, 0.2)),
	   }

	   // Create a scene, containing a camera and a list of objects to render.
	   scn := scene{cam(v(6.0, 2.0, 2.5), v(0.0, 1.0, -1.0), 60.0, 0.1, 1.0), objList}

	*/
	// Get the current time, use this to get the elapsed time later.
	startTimeGo := time.Now()

	img := render(randScene())

	// Print how long it took to raycast.
	elapsedGo := time.Since(startTimeGo)
	fmt.Println("Time spent raycasting:", elapsedGo.Seconds(), "s")

	// Flip the image because we want 0, 0 to be the bottom left.
	img = imaging.FlipV(img)

	// Save the file to the destination given in the argument.
	err = saveFile(os.Args[1], img)
	check(err)
}

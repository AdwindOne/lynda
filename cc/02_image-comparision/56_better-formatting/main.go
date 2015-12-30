package main

import (
	"fmt"
	stdimage "image"
	"image/jpeg"
	"log"
	"os"
	"time"
	"strings"
	"path/filepath"
	"sync"
)

// at terminal:
// go run -race main.go

type pixel struct {
	r, g, b, a uint32
}

type image struct {
	name   string
	pixels []pixel
	width  int
	height int
}

type result struct {
	needle, haystack *image
	hx int
	diff float64
}

var wg sync.WaitGroup
const dir = "../00_images/"
const threshold = 2000

func main() {
	start := time.Now()

	ch := getImages()
	images := xImages(ch)

	for _, v := range images {
		fmt.Println(*v)
		fmt.Println("\n\n\n\n\n\n\n")
		fmt.Println("***********************")
		fmt.Println("\n\n\n\n\n\n\n")
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func getImages() chan *image {

	paths, err := getPaths()
	if err != nil {
		log.Println("Error getting paths", err)
	}

	ch := make(chan *image)

	for _, path := range paths {
		wg.Add(1)
		// worker
		go func(path string) {
			defer wg.Done()
			ch <- getPixels(path)
		}(path)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	/*
	"The structure of the code above is a common and idiomatic pattern
	for looping in parallel when we don't know the number of iterations."
	source: Donovan & Kernighan, "The Go Programming Language", p.238
	*/

	return ch
}

func getPaths() ([]string, error) {
	var paths []string

	wf := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	}

	if err := filepath.Walk(dir, wf); err != nil {
		return nil, err
	}

	return paths, nil
}

func getPixels(path string) *image {
	img := loadImage(path)
	bounds := img.Bounds()
	fmt.Println(bounds.Dx(), " x ", bounds.Dy()) // debugging
	pixels := make([]pixel, bounds.Dx()*bounds.Dy())

	for i := 0; i < bounds.Dx()*bounds.Dy(); i++ {
		x := i % bounds.Dx()
		y := i / bounds.Dx()
		r, g, b, a := img.At(x, y).RGBA()
		pixels[i].r = r
		pixels[i].g = g
		pixels[i].b = b
		pixels[i].a = a
	}

	xs := strings.Split(path, "/")
	name := xs[(len(xs) - 1)]
	image := image{
		name:   name,
		pixels: pixels,
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}

	return &image
}

func loadImage(filename string) stdimage.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func xImages(ch chan *image) []*image {

	var images []*image

	for imgPtr := range ch {
		images = append(images, imgPtr)
	}

	return images
}
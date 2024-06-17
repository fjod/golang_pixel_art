package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"sync"
)

func main() {
	window, img := parseAndValidateInput()

	bounds := (*img).Bounds()
	emptyPixelated := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	fmt.Println("created new image with bounds ", emptyPixelated.Bounds())

	rgba := imageToRGBA(*img)
	source := NewImgAsEnumerable(rgba, *window)
	dest := NewImgAsEnumerable(emptyPixelated, *window)
	if source.img == nil || dest.img == nil {
		fmt.Printf("source or dest image is nil")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	ch := make(chan TwoParts)
	for step(source) {
		sourcePart := source.Value()
		destPart := dest.Value()
		parts := TwoParts{Src: sourcePart, Dest: destPart}
		wg.Add(1)
		go send(ch, parts) // async processing all image parts
		go receive(ch, &wg)
	}

	wg.Wait() // wait for all parts to be processed
	close(ch)

	writeOutput(emptyPixelated)
}

func step(source *ImgAsEnumerable) bool {
	c, _ := source.Next()
	return c
}

func send(ch chan<- TwoParts, v TwoParts) {
	ch <- v
}

func receive(ch chan TwoParts, wg *sync.WaitGroup) {
	val := <-ch
	Process(&val.Src, &val.Dest)
	wg.Done()
}

func writeOutput(emptyPixelated *image.RGBA) {
	_, e := os.Stat("./output.jpg")
	if e != nil {
		if os.IsExist(e) {
			err := os.Remove("./output.jpg")
			if err != nil {
				fmt.Printf("old img cant be deleted")
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	out, err := os.Create("./output.jpg")
	if err != nil {
		fmt.Printf("Cant create output file")
		fmt.Println(err)
		os.Exit(1)
	}

	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(out, emptyPixelated, &opt) // put quality to 80%
	if err != nil {
		fmt.Printf("Cant save output image to file")
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseAndValidateInput() (*int, *image.Image) {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: -i [path-to-image] -w [optional output pixel size]\n")
		os.Exit(1)
	}

	inputFile := flag.String("i", "pixel.jpg", "Enter input file")
	window := flag.Int("w", 10, "Enter pixel size")
	flag.Parse()

	fmt.Printf("window is %d\n", *window)

	file, err := os.Open(*inputFile) // For read access.
	if err != nil {
		fmt.Printf("cant open file %s\n", os.Args[1])
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	img, _, imgError := image.Decode(file)
	if imgError != nil {
		fmt.Printf("img cant be decoded")
		fmt.Println(imgError)
		os.Exit(1)
	}
	return window, &img
}

func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

const (
	light = "Ñ@#W$9876543210?!abc;:+=-,._    "
)

func main() {
	var size int
	var fileimage string
	flag.IntVar(&size, "s", 0, "size of the result. This value should not be less then the image size. (default it the image size)")
	flag.StringVar(&fileimage, "image", "", "image file")

	flag.Parse()

	if len(fileimage) == 0 {
		log.Fatalf("missing image")
	}

	f, err := os.Open(fileimage)

	if err != nil {
		log.Fatalf("could not open the image: %v", err)
	}

	defer f.Close()

	c, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Fatalf("could not decode config the image: %v", err)
	}
	h, w := c.Height, c.Width

	if size == 0 {
		size = w
	}

	if size > w {
		log.Fatal("this image is too small")
	}

	f.Seek(0, 0)

	img, _, err := image.Decode(f)

	if err != nil {
		log.Fatalf("could not decode the image: %v", err)
	}

	lightsize := float64(len(light) - 1)

	mw, mh := w/size, h/size

	var line string
	for y := 0; y < h; y += mh {
		line = ""
		for x := 0; x < w; x += mw {
			avg := average(img, x, y, mw, mh)
			pos := mapRange(avg, 0, 255, lightsize, 0)
			pos = math.Abs(pos)

			line += string(light[int(pos)])
		}
		fmt.Fprintf(os.Stdout, "%s\n", line)
	}
}

func average(img image.Image, x, y, mw, mh int) float64 {
	var total int
	for i := 0; i < mw; i++ {
		for j := 0; j < mw; j++ {
			clr := img.At(x+i, y+j)
			r, g, b, _ := clr.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			avg := (r + g + b) / 3
			total += int(avg)
		}
	}
	return float64(total) / float64(mw*mh)
}

func mapRange(s, a1, a2, b1, b2 float64) float64 {
	return b1 + (s-a1)*(b2-b1)/(a2-a1)
}

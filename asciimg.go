package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

const light = "Ñ@#W$9876543210?!abc;:+=-,._    "

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: asciimg <image>")
	}
	fn := os.Args[1]

	f, err := os.Open(fn)

	if err != nil {
		log.Fatalf("could not open the image: %v", err)
	}

	defer f.Close()

	c, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Fatalf("could not decode config the image: %v", err)
	}

	f.Seek(0, 0)

	img, _, err := image.Decode(f)

	if err != nil {
		log.Fatalf("could not decode the image: %v", err)
	}

	size := float64(len(light) - 1)

	var line string
	for y := 0; y < c.Height; y++ {
		line = ""
		for x := 0; x < c.Width; x++ {
			clr := img.At(x, y)
			r, g, b, _ := clr.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			avg := (r + g + b) / 3

			pos := mapRange(float64(avg), 0, 255, size, 0)
			pos = math.Abs(pos)
			line += string(light[int(pos)])
		}
		fmt.Fprintf(os.Stdout, "%s\n", line)
	}

}

func mapRange(s, a1, a2, b1, b2 float64) float64 {
	return b1 + (s-a1)*(b2-b1)/(a2-a1)
}

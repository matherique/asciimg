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

const (
	light              = "Ñ@#W$9876543210?!abc;:+=-,._    "
	MAX_TO_OPTIMIZE_2X = 100
	MAX_TO_OPTIMIZE_3X = 200
	MAX_TO_OPTIMIZE_4X = 500
	MAX_TO_OPTIMIZE_5X = 1000
)

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

	h, w := c.Height, c.Width

	opt := calculateOptimizeIndex(w, h)

	fmt.Println(opt)
	var line string
	for y := 0; y < h; y += opt {
		line = ""
		for x := 0; x < w; x += opt {
			avg := average(img, opt, x, y)
			pos := mapRange(float64(avg), 0, 255, size, 0)
			pos = math.Abs(pos)

			line += string(light[int(pos)])
		}
		fmt.Fprintf(os.Stdout, "%s\n", line)
	}
}

func average(img image.Image, opt, x, y int) float64 {
	m := 1
	if opt > 1 {
		m = opt - 2
	}

	var total int
	for i := 0; i <= m; i++ {
		for j := 0; j <= m; j++ {
			clr := img.At(x+i, y+j)
			r, g, b, _ := clr.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			avg := (r + g + b) / 3
			total += int(avg)
		}
	}
	return float64(total) / 9
}

func calculateOptimizeIndex(w, h int) int {
	switch true {
	case w > MAX_TO_OPTIMIZE_5X:
		return 5
	case w > MAX_TO_OPTIMIZE_4X:
		return 4
	case w > MAX_TO_OPTIMIZE_3X:
		return 3
	case w > MAX_TO_OPTIMIZE_2X:
		return 2
	}

	return 1
}

func mapRange(s, a1, a2, b1, b2 float64) float64 {
	return b1 + (s-a1)*(b2-b1)/(a2-a1)
}

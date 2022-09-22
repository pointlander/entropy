// Copyright 2022 The Entropy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	input, err := os.Open("lenna.png")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	img, err := png.Decode(input)
	if err != nil {
		panic(err)
	}

	b := img.Bounds()
	set := image.NewRGBA(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			original := img.At(x, y)
			pixel := color.GrayModel.Convert(original)
			set.Set(x, y, pixel)
		}
	}

	output, err := os.Create("gray.jpg")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	png.Encode(output, set)
}

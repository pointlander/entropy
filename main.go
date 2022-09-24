// Copyright 2022 The Entropy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"math/rand"
	"os"
	"sort"

	"github.com/mjibson/go-dsp/fft"
)

func main() {
	os.MkdirAll("gray", 0755)
	Process("image01")
	Process("image02")
	Process("image03")
	Process("image04")
}

// Entropy returns the entropy of the given image
func Entropy(x [][]complex128) complex128 {
	var sum complex128
	for _, r := range x {
		for _, c := range r {
			sum += c * cmplx.Log(c)
		}
	}
	return -sum
}

func Process(name string) {
	fmt.Println(name)
	rnd := rand.New(rand.NewSource(1))
	input, err := os.Open(fmt.Sprintf("images/%s.png", name))
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
	xx := make([][]complex128, b.Max.Y)
	zz := make([]uint8, b.Max.Y*b.Max.X)
	for y := 0; y < b.Max.Y; y++ {
		xx[y] = make([]complex128, b.Max.X)
		for x := 0; x < b.Max.X; x++ {
			original := img.At(x, y)
			pixel := color.GrayModel.Convert(original)
			gray, _ := pixel.(color.Gray)
			xx[y][x] = complex(float64(gray.Y), 0)
			zz[y*b.Max.X+x] = gray.Y
			set.Set(x, y, pixel)
		}
	}

	y := fft.FFT2(xx)
	entropy := Entropy(y)
	fmt.Println("original", entropy, cmplx.Abs(entropy), cmplx.Phase(entropy))

	rnd.Shuffle(len(zz), func(i, j int) {
		zz[i], zz[j] = zz[j], zz[i]
	})
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			xx[y][x] = complex(float64(zz[y*b.Max.X+x]), 0)
		}
	}
	y = fft.FFT2(xx)
	entropy = Entropy(y)
	fmt.Println("shuffled", entropy, cmplx.Abs(entropy), cmplx.Phase(entropy))

	sort.Slice(zz, func(i, j int) bool {
		return zz[i] < zz[j]
	})
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			xx[y][x] = complex(float64(zz[y*b.Max.X+x]), 0)
		}
	}
	y = fft.FFT2(xx)
	entropy = Entropy(y)
	fmt.Println("sorted", entropy, cmplx.Abs(entropy), cmplx.Phase(entropy))

	output, err := os.Create(fmt.Sprintf("gray/%s.png", name))
	if err != nil {
		panic(err)
	}
	defer output.Close()
	png.Encode(output, set)
}

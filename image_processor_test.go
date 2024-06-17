package main

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("SumOfEmptySlice", func(t *testing.T) {
		var input []uint8
		expected := uint32(0)
		actual := sum(&input)
		if actual != expected {
			t.Errorf("sum(%v) = %d; expected %d", input, actual, expected)
		}
	})

	t.Run("SumOfSingleElement", func(t *testing.T) {
		input := []uint8{42}
		expected := uint32(42)
		actual := sum(&input)
		if actual != expected {
			t.Errorf("sum(%v) = %d; expected %d", input, actual, expected)
		}
	})

	t.Run("SumOfMultipleElements", func(t *testing.T) {
		input := []uint8{1, 2, 3, 4, 5, math.MaxUint8 / 3, math.MaxUint8 / 2, math.MaxUint8 / 2}
		expected := uint32(15 + math.MaxUint8/3 + math.MaxUint8/2 + math.MaxUint8/2)
		actual := sum(&input)
		if actual != expected {
			t.Errorf("sum(%v) = %d; expected %d", input, actual, expected)
		}
	})

	t.Run("SumOfMaxUint32Values", func(t *testing.T) {
		input := []uint8{math.MaxUint8, math.MaxUint8}
		expected := uint32(math.MaxUint8) * 2
		actual := sum(&input)
		if actual != expected {
			t.Errorf("sum(%v) = %d; expected %d", input, actual, expected)
		}
	})
}

func TestCalcMedianColor(t *testing.T) {
	t.Run("EmptyImgPart", func(t *testing.T) {
		input := &ImgPart{}
		expected := color.RGBA{}
		actual := calcMedianColor(input)
		if actual != expected {
			t.Errorf("calcMedianColor(%v) = %v; expected %v", input, actual, expected)
		}
	})

	t.Run("SinglePixelImgPart", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.RGBA{R: 255, A: 255})
		input := &ImgPart{
			Img: img,
		}
		expected := color.RGBA{R: 255, A: 255}
		actual := calcMedianColor(input)
		if actual != expected {
			t.Errorf("calcMedianColor(%v) = %v; expected %v", input, actual, expected)
		}
	})

	t.Run("ManyPixels", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				img.Set(x, y, color.RGBA{R: 127, G: 127, B: 127, A: 255})
			}
		}

		input := &ImgPart{
			Img: img,
		}
		expected := color.RGBA{R: 127, G: 127, B: 127, A: 255}
		actual := calcMedianColor(input)
		if actual != expected {
			t.Errorf("calcMedianColor(%v) = %v; expected %v", input, actual, expected)
		}
	})
}

package main

import (
	"image/color"
)

func Process(input *ImgPart, output *ImgPart) {
	medianColor := calcMedianColor(input)
	for y := output.Img.Rect.Min.Y; y < output.Img.Rect.Max.Y; y++ {
		for x := output.Img.Rect.Min.X; x < output.Img.Rect.Max.X; x++ {
			output.Img.Set(x, y, medianColor)
		}
	}
}

func calcMedianColor(input *ImgPart) color.Color {
	if input.Img == nil {
		return color.RGBA{}
	}

	pixelCount := (input.Img.Rect.Max.Y - input.Img.Rect.Min.Y) * (input.Img.Rect.Max.X - input.Img.Rect.Min.X)
	var totalR = make([]uint8, 0, pixelCount)
	var totalG = make([]uint8, 0, pixelCount)
	var totalB = make([]uint8, 0, pixelCount)
	var totalA = make([]uint8, 0, pixelCount)

	bounds := input.Img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			current := input.Img.At(x, y)
			r, g, b, a := current.RGBA()

			// https://stackoverflow.com/questions/33186783/get-a-pixel-array-from-from-golang-image-image
			totalR = append(totalR, uint8(r/257))
			totalG = append(totalG, uint8(g/257))
			totalB = append(totalB, uint8(b/257))
			totalA = append(totalA, uint8(a/257))
		}
	}
	return color.RGBA{
		R: uint8(sum(&totalR) / uint32(len(totalR))),
		G: uint8(sum(&totalG) / uint32(len(totalG))),
		B: uint8(sum(&totalB) / uint32(len(totalB))),
		A: uint8(sum(&totalA) / uint32(len(totalA))),
	}
}

func sum(input *[]uint8) uint32 {
	var total uint32
	for _, v := range *input {
		total += uint32(v)
	}
	return total
}

package main

import (
	"image"
)

type ImgAsEnumerable struct {
	img      *image.RGBA
	window   int
	rect     image.Rectangle
	xShift   int
	yShift   int
	finished bool
}

type ImgPart struct {
	Img *image.RGBA
}

type TwoParts struct {
	Src  ImgPart
	Dest ImgPart
}

func NewImgAsEnumerable(input *image.RGBA, window int) *ImgAsEnumerable {
	if input == nil {
		return &ImgAsEnumerable{}
	}
	if window < 1 {
		return &ImgAsEnumerable{}
	}
	p := new(ImgAsEnumerable)
	p.img = input
	p.window = window
	b := input.Bounds()
	p.rect = image.Rect(0, 0, b.Max.X, b.Max.Y)
	return p
}

func (Img *ImgAsEnumerable) Next() (bool, error) {
	if Img.finished {
		return false, nil
	}
	if (Img.yShift) >= Img.rect.Dy() {
		Img.finished = true
		return false, nil
	}
	return true, nil
}

func (Img *ImgAsEnumerable) Value() ImgPart {
	right := Img.xShift + Img.window
	if Img.xShift+Img.window > Img.rect.Dx() {
		right = Img.rect.Dx()
	}
	bot := Img.yShift + Img.window
	if Img.yShift+Img.window > Img.rect.Dy() {
		bot = Img.rect.Dy()
	}
	current := Img.img.SubImage(image.Rect(Img.xShift, Img.yShift, right, bot)).(*image.RGBA)
	Img.xShift += Img.window
	if Img.xShift >= Img.rect.Dx() {
		Img.xShift = 0
		Img.yShift += Img.window
		if Img.yShift >= Img.rect.Dy() {
			Img.yShift = 0
			Img.finished = true
		}
	}

	var ret ImgPart
	ret.Img = current
	return ret
}

func (Img *ImgAsEnumerable) Reset() {
	Img.xShift = 0
	Img.yShift = 0
	Img.finished = false
}

package main

import (
	"image"
	"reflect"
	"testing"
)

var SampleRect = image.NewRGBA(image.Rect(0, 0, 10, 10))
var SampleNonRect = image.NewRGBA(image.Rect(0, 0, 15, 10))

func TestCtor(t *testing.T) {
	tests := []struct {
		name     string
		input    *image.RGBA
		window   int
		expected *ImgAsEnumerable
	}{
		{
			name:     "NilInput",
			input:    nil,
			window:   10,
			expected: &ImgAsEnumerable{},
		},
		{
			name:     "InvalidWindow",
			input:    SampleRect,
			window:   0,
			expected: &ImgAsEnumerable{},
		},
		{
			name:     "ValidInput",
			input:    SampleRect,
			window:   10,
			expected: &ImgAsEnumerable{img: SampleRect, window: 10, rect: image.Rect(0, 0, 100, 100)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewImgAsEnumerable(test.input, test.window)
			if result.img != test.expected.img {
				t.Errorf("NewImgAsEnumerable(%v, %d) = %v, want %v", test.input, test.window, result, test.expected)
			}
			if result.window != test.expected.window {
				t.Errorf("NewImgAsEnumerable(%v, %d) = %v, want %v", test.input, test.window, result, test.expected)
			}
		})
	}
}

func TestNextSanityCheck(t *testing.T) {
	tests := []struct {
		name           string
		input          ImgAsEnumerable
		expectedResult bool
		expectedError  error
	}{
		{
			name: "FinishedEnumeration",
			input: ImgAsEnumerable{
				finished: true,
				img:      SampleRect,
				window:   10,
				rect:     image.Rect(0, 0, 10, 10),
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "YShiftExceedsRect",
			input: ImgAsEnumerable{
				finished: false,
				img:      SampleRect,
				window:   10,
				rect:     image.Rect(0, 0, 10, 10),
				yShift:   10,
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "ValidEnumeration",
			input: ImgAsEnumerable{
				finished: false,
				img:      SampleRect,
				window:   5,
				rect:     image.Rect(0, 0, 10, 10),
			},
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.input.Next()
			if result != test.expectedResult {
				t.Errorf("Next() = %v, want %v", result, test.expectedResult)
			}
			if err != nil {
				t.Errorf("Next() error = %v, want %v", err, test.expectedError)
			}
		})
	}
}

func TestNextValidEnumeration(t *testing.T) {
	tests := []struct {
		name     string
		input    *image.RGBA
		window   int
		expected int // amount of enumerations
	}{
		{
			name:     "ValidEnumerationOnEvenWindow",
			input:    SampleRect,
			window:   5,
			expected: 4,
		},
		{
			name:     "ValidEnumerationOnOddWindow",
			input:    SampleRect,
			window:   3,
			expected: 16,
		},
		{
			name:     "ValidEnumerationOnEvenWindow",
			input:    SampleNonRect,
			window:   5,
			expected: 6,
		},
		{
			name:     "ValidEnumerationOnOddWindow",
			input:    SampleNonRect,
			window:   6,
			expected: 6,
		},
	}

	for _, test := range tests {
		var e error
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Window is %d, image is %d", test.window, test.input.Rect.Dx())
			result := NewImgAsEnumerable(test.input, test.window)
			modX := test.input.Rect.Dx() % test.window
			modY := test.input.Rect.Dy() % test.window
			validRects := []struct {
				x int
				y int
			}{
				{test.window, test.window},
				{test.window, modY},
				{modX, test.window},
				{modX, modY},
			}
			count := 0
			for c, _ := true, e; c; c, _ = result.Next() {
				v := result.Value()
				xSize := v.Img.Rect.Dx()
				ySize := v.Img.Rect.Dy()
				t.Logf("The square is x = %d y = %d", xSize, ySize)
				foundValid := false
				for _, rect := range validRects {
					if xSize == rect.x && ySize == rect.y {
						foundValid = true
					}
				}
				if !foundValid {
					t.Errorf("Incorrect result rect")
				}
				count += 1
			}

			if count != test.expected {
				t.Errorf("Next() = %v, want %v", count, test.expected)
			}

			if result.finished != true {
				t.Errorf("Incorrect result.finished")
			}
		})
	}
}

func TestResetImgAsEnumerable(t *testing.T) {
	tests := []struct {
		name     string
		input    ImgAsEnumerable
		expected ImgAsEnumerable
	}{
		{
			name: "ResetFromFinished",
			input: ImgAsEnumerable{
				finished: true,
				xShift:   5,
				yShift:   5,
			},
			expected: ImgAsEnumerable{
				finished: false,
				xShift:   0,
				yShift:   0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.input.Reset()
			if !reflect.DeepEqual(test.input, test.expected) {
				t.Errorf("Reset() = %v, want %v", test.input, test.expected)
			}
		})
	}
}

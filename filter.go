package gogh

import (
	"github.com/ironpark/gogh/mask"
	"math"
)

const (
	BLUR_BOX = 0
)

func (src *Img) Filter(mask interface{}, value ...int) *Img {
	switch v := mask.(type) {
	case [][]float32:
		//convolution mask
		return convolution(src, v)
	case int, string:
		switch mask {
		case "sobel":

		case "Gaussian":

		case "box":
			return src.Blur(value[0], BLUR_BOX)
		}
	default:
		//err
	}

	return src
}

func (src *Img) convolution1mask(x, y, center int, kernel [][]float32) float32 {
	//float32 := 0
	bounds := src.Bounds()
	var cPixel float32
	for ky, Es := range kernel {
		for kx, E := range Es {
			xe1 := kx - center
			ye2 := ky - center

			selx := x + xe1
			sely := y + ye2

			if selx < 0 {
				selx = 0
			}
			if selx > (bounds.Max.X - 1) {
				selx = (bounds.Max.X - 1)
			}
			if sely < 0 {
				sely = 0
			}
			if sely > (bounds.Max.Y - 1) {
				sely = (bounds.Max.Y - 1)
			}
			gray := src.At(selx, sely).Gray()
			cPixel += float32(gray) * E
		}
	}
	return cPixel
}
func convolution(src *Img, kernel [][]float32) *Img {
	size := len(kernel)
	//kernel center point
	kernelCP := int(float32(size) / float32(2))
	normalization := false
	normal := float32(0)
	for _, Es := range kernel {
		for _, E := range Es {
			normal += E
		}
	}

	if normal > 1 {
		normalization = true
	}

	v := clone(src)
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Convolution
			cPixel := src.convolution1mask(x, y, kernelCP, kernel)
			//normalization
			if normalization {
				cPixel = cPixel / normal
			}
			c := int(math.Abs(float64(cPixel)))
			if 255 < c {
				c = 255
			}
			//set pixel

			v.At(x, y).Set(c, c, c)
		}
	}
	return src
}

func (src *Img) FindEdge(maskX, maskY [][]float32, t int) *Img {
	bounds := src.Bounds()
	edgeImage := NewImg(bounds)
	kernelCP := int(float32(len(maskX)) / float32(2))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Convolution
			dx := src.convolution1mask(x, y, kernelCP, maskX)
			dy := src.convolution1mask(x, y, kernelCP, maskY)
			//강도
			mag := int(math.Abs(float64(dx + dy)))
			//set pixel
			if mag >= t {
				edgeImage.At(x, y).Set(255, 255, 255)
			} else {
				edgeImage.At(x, y).Set(0, 0, 0)
			}
		}
	}
	src = edgeImage
	return edgeImage
}
func (src *Img) Blur(size, blurtype int) *Img {
	switch blurtype {
	case BLUR_BOX:
		return src.Filter(mask.GenBoxBlurMask(size))
	}
	//unknown filter
	return nil
}

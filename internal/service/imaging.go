package service

import (
	"image"

	"github.com/disintegration/imaging"
)

func Blur(img image.Image, sigma float64) *image.NRGBA {
	return imaging.Blur(img, sigma)
}

func Sharpen(img image.Image, sigma float64) *image.NRGBA {
	return imaging.Sharpen(img, sigma)
}

func Resize(img image.Image, width int, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Linear)
}

func Crop(img image.Image, width int, height int, anchor string) *image.NRGBA {
	var an imaging.Anchor

	switch anchor {
	case "center":
		an = imaging.Center
	case "top":
		an = imaging.Top
	case "bottom":
		an = imaging.Bottom
	case "top-left":
		an = imaging.TopLeft
	case "top-right":
		an = imaging.TopRight
	case "bottom-left":
		an = imaging.BottomLeft
	case "bottom-right":
		an = imaging.BottomRight
	}

	return imaging.CropAnchor(img, width, height, an)
}

func Grayscale(img image.Image) *image.NRGBA {
	return imaging.Grayscale(img)
}

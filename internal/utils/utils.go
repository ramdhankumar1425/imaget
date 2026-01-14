package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramdhankumar1425/imaget/internal/model"
)

func GetJobParams(c *gin.Context, jobType string) (model.TransformParams, error) {
	var params model.TransformParams

	switch jobType {
	case "blur":
		sigmaStr := c.PostForm("sigma")
		if sigmaStr == "" {
			return nil, fmt.Errorf("sigma required")
		}

		sigma, err := strconv.ParseFloat(sigmaStr, 64)
		if err != nil || sigma <= 0 {
			return nil, fmt.Errorf("invalid sigma")
		}

		params = model.BlurParams{Sigma: sigma}
	case "sharpen":
		sigmaStr := c.PostForm("sigma")
		if sigmaStr == "" {
			return nil, fmt.Errorf("sigma required")
		}

		sigma, err := strconv.ParseFloat(sigmaStr, 64)
		if err != nil || sigma <= 0 {
			return nil, fmt.Errorf("invalid sigma")
		}

		params = model.SharpenParams{Sigma: sigma}
	case "resize":
		wStr := c.PostForm("width")
		hStr := c.PostForm("height")

		w, wErr := strconv.Atoi(wStr)
		h, hErr := strconv.Atoi(hStr)
		if wErr != nil || hErr != nil || w <= 0 || h <= 0 {
			return nil, fmt.Errorf("invalid width or height")
		}

		params = model.ResizeParams{
			Width:  w,
			Height: h,
		}
	case "crop":
		wStr := c.PostForm("width")
		hStr := c.PostForm("height")
		anchor := c.DefaultPostForm("anchor", "center")

		w, wErr := strconv.Atoi(wStr)
		h, hErr := strconv.Atoi(hStr)
		if wErr != nil || hErr != nil || w <= 0 || h <= 0 {
			return nil, fmt.Errorf("invalid width or height")
		}

		validAnchors := map[string]bool{
			"center":       true,
			"top":          true,
			"bottom":       true,
			"left":         true,
			"right":        true,
			"top-left":     true,
			"top-right":    true,
			"bottom-left":  true,
			"bottom-right": true,
		}
		if !validAnchors[anchor] {
			return nil, fmt.Errorf("invalid anchor")
		}

		params = model.CropParams{
			Width:  w,
			Height: h,
			Anchor: anchor,
		}
	case "grayscale":
		params = model.GrayscaleParams{}
	default:
		return nil, fmt.Errorf("invalid transformation")
	}

	return params, nil
}

func GetFileExtension(format string) (string, error) {
	switch format {
	case "jpeg", "jpg":
		return "jpg", nil
	case "png":
		return "png", nil
	default:
		return "", fmt.Errorf("invalid file extension")
	}
}

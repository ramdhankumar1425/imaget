package utils

import (
	"context"
	"fmt"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

func DownloadAndDecodeImage(ctx context.Context, url string) (image.Image, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download failed: %s", resp.Status)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

func IsValidImageKitRawURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	prefix := os.Getenv("IMAGEKIT_URL_PREFIX")
	if prefix == "" {
		log.Fatal("IMAGEKIT_URL_PREFIX env not set")
	}

	return strings.HasPrefix(u.String(), prefix)
}

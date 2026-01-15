package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"time"

	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/packages/param"
	"github.com/ramdhankumar1425/imaget/internal/infra"
	"github.com/ramdhankumar1425/imaget/internal/model"
	"github.com/ramdhankumar1425/imaget/internal/service"
	"github.com/ramdhankumar1425/imaget/internal/utils"
)

func Worker() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("worker panic:", r)
			go Worker()
		}
	}()

	for job := range Jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

		func() {
			defer cancel()

			log.Printf("Starting Job with ID: %v", job.ID)
			infra.RDB.HSet(ctx, "job:"+job.ID,
				map[string]any{
					"status":    "processing",
					"startedAt": time.Now().Unix(),
				},
			)

			// download the file
			img, format, err := utils.DownloadAndDecodeImage(ctx, job.RawURL)
			if err != nil {
				log.Println("cannot download/decode image")
				infra.RDB.HSet(ctx, "job:"+job.ID,
					map[string]any{
						"status": "failed",
						"error":  "cannot download/decode image",
					},
				)
				return
			}

			fileExt, err := utils.GetFileExtension(format)
			if err != nil {
				log.Println("invalid file extension")
				infra.RDB.HSet(ctx, "job:"+job.ID,
					map[string]any{
						"status": "failed",
						"error":  "invalid file extension",
					},
				)
				return
			}

			var result *image.NRGBA
			switch p := job.Params.(type) {
			case model.BlurParams:
				result = service.Blur(img, p.Sigma)
			case model.SharpenParams:
				result = service.Sharpen(img, p.Sigma)
			case model.ResizeParams:
				result = service.Resize(img, p.Width, p.Height)
			case model.CropParams:
				result = service.Crop(img, p.Width, p.Height, p.Anchor)
			case model.GrayscaleParams:
				result = service.Grayscale(img)
			}

			if result == nil {
				infra.RDB.HSet(ctx, "job:"+job.ID,
					map[string]any{
						"status": "failed",
						"error":  "invalid job type",
					},
				)
				return
			}

			var buf bytes.Buffer
			var fileTypeErr error
			switch format {
			case "jpeg", "jpg":
				fileTypeErr = jpeg.Encode(&buf, result, &jpeg.Options{Quality: 80})
			case "png":
				fileTypeErr = png.Encode(&buf, result)
			default:
				fileTypeErr = fmt.Errorf("unsupported file type %s", format)
			}

			if fileTypeErr != nil {
				infra.RDB.HSet(
					ctx,
					"job:"+job.ID,
					map[string]any{
						"status": "failed",
						"error":  fileTypeErr.Error(),
					},
				)
				return
			}

			reader := bytes.NewReader(buf.Bytes())

			res, err := infra.ImageKit.Files.Upload(ctx, imagekit.FileUploadParams{
				File:     reader,
				FileName: job.ID + "." + fileExt,
				Folder: param.Opt[string]{
					Value: "/result",
				},
			})
			if err != nil {
				infra.RDB.HSet(ctx, "job:"+job.ID,
					map[string]any{
						"status": "failed",
						"error":  err.Error(),
					},
				)
				return
			}

			infra.RDB.HSet(
				ctx,
				"job:"+job.ID,
				map[string]any{
					"status":      "completed",
					"resultUrl":   res.URL,
					"completedAt": time.Now().Unix(),
				},
			)

			log.Printf("Completed Job with ID: %v", job.ID)
		}()
	}
}

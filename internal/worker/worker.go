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
)

func Worker() {
	ctx := context.Background()

	for job := range Jobs {
		log.Printf("Starting Job with ID: %v", job.Meta.ID)
		infra.RDB.HSet(ctx, "job:"+job.Meta.ID, "status", "processing", "startedAt", time.Now().Unix())

		var result *image.NRGBA
		switch p := job.Params.(type) {
		case model.BlurParams:
			result = service.Blur(job.Image, p.Sigma)
		case model.SharpenParams:
			result = service.Sharpen(job.Image, p.Sigma)
		case model.ResizeParams:
			result = service.Resize(job.Image, p.Width, p.Height)
		case model.CropParams:
			result = service.Crop(job.Image, p.Width, p.Height, p.Anchor)
		case model.GrayscaleParams:
			result = service.Grayscale(job.Image)
		}

		if result == nil {
			infra.RDB.HSet(ctx, "job:"+job.Meta.ID, "status", "failed", "error", "invalid job type")
			continue
		}

		var buf bytes.Buffer
		var err error
		switch job.Meta.FileType {
		case "jpeg", "jpg":
			err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: 90})
		case "png":
			err = png.Encode(&buf, result)
		default:
			err = fmt.Errorf("unsupported file type %s", job.Meta.FileType)

		}
		if err != nil {
			infra.RDB.HSet(
				ctx,
				"job:"+job.Meta.ID,
				"status", "failed",
				"error", err.Error(),
			)
			continue
		}

		reader := bytes.NewReader(buf.Bytes())

		res, err := infra.ImageKit.Files.Upload(ctx, imagekit.FileUploadParams{
			File:     reader,
			FileName: job.Meta.ID + "." + job.Meta.FileExt,
			Folder: param.Opt[string]{
				Value: "/result",
			},
		})
		if err != nil {
			infra.RDB.HSet(ctx, "job:"+job.Meta.ID, "status", "failed", "error", err.Error())
			continue
		}

		infra.RDB.HSet(
			ctx,
			"job:"+job.Meta.ID,
			"status", "completed",
			"resultUrl", res.URL,
			"completedAt", time.Now().Unix(),
		)

		log.Printf("Completed Job with ID: %v", job.Meta.ID)
	}
}

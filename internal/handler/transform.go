package handler

import (
	"image"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ramdhankumar1425/imaget/internal/infra"
	"github.com/ramdhankumar1425/imaget/internal/model"
	"github.com/ramdhankumar1425/imaget/internal/utils"
	"github.com/ramdhankumar1425/imaget/internal/worker"

	_ "image/jpeg"
	_ "image/png"
)

func HandleTransform(c *gin.Context) {
	log.Println("Handling job...")
	ctx := c.Request.Context()

	jobType := c.PostForm("type")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file not provided"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot open file"})
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid image"})
		return
	}

	jobParams, err := utils.GetJobParams(c, jobType)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	fileExt, err := utils.GetFileExtension(format)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	jobMeta := model.JobMeta{
		ID:          uuid.NewString(),
		Status:      "pending",
		Type:        jobType,
		FileSizeMB:  (float64(fileHeader.Size) / 1024) / 1024,
		FileType:    format,
		FileExt:     fileExt,
		ResultURL:   "",
		Error:       "",
		QueuedAt:    time.Now().Unix(),
		StartedAt:   0,
		CompletedAt: 0,
	}

	job := model.Job{
		Meta:   jobMeta,
		Image:  img,
		Params: jobParams,
	}

	infra.RDB.HSet(ctx, "job:"+jobMeta.ID, map[string]any{
		"id":          jobMeta.ID,
		"status":      jobMeta.Status,
		"type":        jobMeta.Type,
		"fileSizeMb":  jobMeta.FileSizeMB,
		"fileType":    jobMeta.FileType,
		"fileExt":     fileExt,
		"resultUrl":   jobMeta.ResultURL,
		"error":       jobMeta.Error,
		"queuedAt":    jobMeta.QueuedAt,
		"startedAt":   jobMeta.StartedAt,
		"completedAt": jobMeta.CompletedAt,
	})

	worker.Jobs <- job

	log.Println("Job queued, returning response")
	c.JSON(200, gin.H{"jobId": jobMeta.ID})
}

func HandleGetResult(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	data, err := infra.RDB.HGetAll(ctx, "job:"+id).Result()
	if err != nil || len(data) == 0 {
		c.JSON(404, gin.H{"error": "job not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":          data["id"],
		"status":      data["status"],
		"type":        data["type"],
		"fileSizeMb":  data["fileSizeMb"],
		"fileType":    data["fileType"],
		"fileExt":     data["fileExt"],
		"resultUrl":   data["resultUrl"],
		"error":       data["error"],
		"queuedAt":    data["queuedAt"],
		"startedAt":   data["startedAt"],
		"completedAt": data["completedAt"],
	})
}

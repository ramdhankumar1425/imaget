package handler

import (
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
	ctx := c.Request.Context()

	jobType := c.PostForm("type")
	if jobType == "" {
		log.Println("transform type not provided")
		c.JSON(400, gin.H{"error": "transform type not provided"})
		return
	}

	fileUrl := c.PostForm("fileUrl")
	if !utils.IsValidImageKitRawURL(fileUrl) {
		log.Println("invalid file url")
		c.JSON(400, gin.H{"error": "invalid file url"})
		return
	}

	jobParams, err := utils.GetJobParams(c, jobType)
	if err != nil {
		log.Println("invalid job params")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	job := model.Job{
		ID:          uuid.NewString(),
		Status:      "pending",
		Type:        jobType,
		RawURL:      fileUrl,
		ResultURL:   "",
		Error:       "",
		Params:      jobParams,
		QueuedAt:    time.Now().Unix(),
		StartedAt:   0,
		CompletedAt: 0,
	}

	if err := infra.RDB.HSet(ctx, "job:"+job.ID, map[string]any{
		"id":          job.ID,
		"status":      job.Status,
		"type":        job.Type,
		"rawUrl":      job.RawURL,
		"resultUrl":   job.ResultURL,
		"error":       job.Error,
		"queuedAt":    job.QueuedAt,
		"startedAt":   job.StartedAt,
		"completedAt": job.CompletedAt,
	}).Err(); err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	select {
	case worker.Jobs <- job:
		log.Println("Job queued, returning response")
		c.JSON(200, gin.H{"jobId": job.ID})
	default:
		log.Println("server is busy")
		c.JSON(429, gin.H{"error": "server is busy"})
		return
	}
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
		"rawUrl":      data["rawUrl"],
		"resultUrl":   data["resultUrl"],
		"error":       data["error"],
		"queuedAt":    data["queuedAt"],
		"startedAt":   data["startedAt"],
		"completedAt": data["completedAt"],
	})
}

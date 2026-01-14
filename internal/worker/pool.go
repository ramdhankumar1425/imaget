package worker

import (
	"runtime"

	"github.com/ramdhankumar1425/imaget/internal/model"
)

var Jobs chan model.Job

func InitPool() {
	cpus := runtime.NumCPU()

	Jobs = make(chan model.Job, cpus*2)

	for range cpus {
		go Worker()
	}
}

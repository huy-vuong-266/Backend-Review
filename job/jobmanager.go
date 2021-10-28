package job

import (
	"Backend-Review/model"
	"Backend-Review/service"
	"Backend-Review/storage"
)

type JobManager struct{}

func NewJobManager() service.JobManagerInterface {
	return &JobManager{}
}

func (j *JobManager) AddJob(jobs chan<- model.Job, key string, value string) error {

	res := storage.Redis.LPush(key, value)
	job := model.Job{
		Key:   key,
		Value: value,
	}

	jobs <- job

	return res.Err()
}

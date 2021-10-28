package job

import (
	"Backend-Review/constants"
	"Backend-Review/model"
	"Backend-Review/service"
	"encoding/json"
)

type Worker struct{}

func NewWorker() service.WorkerInterface {
	return &Worker{}
}

func (w *Worker) ProcessJob(s *service.Service, jobPool <-chan model.Job) {
	for job := range jobPool {
		switch job.Key {
		case constants.AddFundJobKey:
			reqData := model.FinOrderRequest{}
			json.Unmarshal([]byte(job.Value), &reqData)
			s.FinService.AddFund(reqData.UserID, reqData.Amount)

		case constants.WithdrawJobKey:
			reqData := model.FinOrderRequest{}
			json.Unmarshal([]byte(job.Value), &reqData)
			s.FinService.Withdraw(reqData.UserID, reqData.Amount)
		}
	}
}

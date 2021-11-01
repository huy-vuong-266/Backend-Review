package main

import (
	"Backend-Review/constants"
	fin "Backend-Review/external"
	internal "Backend-Review/internal"
	authen "Backend-Review/internal/authen"
	"Backend-Review/internal/order.go"
	user "Backend-Review/internal/user"
	"Backend-Review/job"
	"Backend-Review/model"
	service "Backend-Review/service"
	storage "Backend-Review/storage"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/go-playground/validator"
	echo "github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	service := &service.Service{
		UserService:   user.NewServiceUser(),
		AuthenService: authen.NewServiceAuthen(),
		FinService:    fin.NewFinService(),
		OrderService:  order.NewServiceOrder(),
		JobManager:    job.NewJobManager(),
		Worker:        job.NewWorker(),
	}

	e.Validator = &model.CustomValidator{
		Validator: validator.New(),
	}

	handler := internal.NewRouter(e, service)

	err := storage.ConnectDB()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("DB connected")

	err = storage.ConnectRedis()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Redis connected")

	server := http.Server{
		Addr:    "127.0.0.1:8003",
		Handler: handler,
	}
	JobPool := make(chan model.Job, 10)
	go func() {
		for {
			ok := storage.Redis.LLen(constants.AddFundJobKey)
			code, _ := ok.Result()
			if code != 0 {
				jobRes := storage.Redis.RPop(constants.AddFundJobKey)
				jobByte, _ := jobRes.Bytes()
				log.Println(string(jobByte))

				select {
				case JobPool <- model.Job{
					Key:   constants.AddFundJobKey,
					Value: string(jobByte),
				}:
				default:
					storage.Redis.LPush(constants.AddFundJobKey, string(jobByte))
				}

			}
		}
	}()
	go func() {
		for {
			ok := storage.Redis.LLen(constants.WithdrawJobKey)
			code, _ := ok.Result()
			if code != 0 {
				jobRes := storage.Redis.RPop(constants.WithdrawJobKey)
				jobByte, _ := jobRes.Bytes()
				log.Println(string(jobByte))

				select {
				case JobPool <- model.Job{
					Key:   constants.WithdrawJobKey,
					Value: string(jobByte),
				}:
				default:
					storage.Redis.LPush(constants.WithdrawJobKey, string(jobByte))
				}
			}

		}
	}()
	go func() {

		for i := 1; i <= runtime.GOMAXPROCS(0); i++ {

			go service.Worker.ProcessJob(service, JobPool)

		}

	}()

	fmt.Println("Server listen and server at port 8003")

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

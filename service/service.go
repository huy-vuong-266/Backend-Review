package service

import "Backend-Review/model"

type UserServiceInterface interface {
	GetUserByPhone(phoneNo string) (interface{}, error)
	GetUserByID(userID string) (interface{}, error)
	CreateUser(user interface{}) error
}

type AuthenServiceInterface interface {
	CreateToken(token interface{}) error
	CheckIfTokenExist(token string) bool
	GetUserIDByAccesstoken(token string) (string, error)
}

type OrderServiceInterface interface {
	CreateOrder(order interface{}) error
	UpdateOrder(order interface{}) error
}

type FinServiceInterface interface {
	AddFund(userID string, amount int64) (int, interface{}, []string)
	Withdraw(userID string, amount int64) (int, interface{}, []string)
}

type WorkerInterface interface {
	ProcessJob(s *Service, jobPool <-chan model.Job)
}

type JobManagerInterface interface {
	AddJob(jobs chan<- model.Job, key string, value string) error
}

type Service struct {
	UserService   UserServiceInterface
	AuthenService AuthenServiceInterface
	FinService    FinServiceInterface
	OrderService  OrderServiceInterface
	JobManager    JobManagerInterface
	Worker        WorkerInterface
}

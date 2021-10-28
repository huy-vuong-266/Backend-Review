package user

import (
	"Backend-Review/model"
	"Backend-Review/service"
	storage "Backend-Review/storage"
)

type ServiceUser struct{}

func NewServiceUser() service.UserServiceInterface {
	return &ServiceUser{}
}

func (s *ServiceUser) GetUserByPhone(phoneno string) (interface{}, error) {
	var user model.User

	err := storage.DB.Where("phone = ?", phoneno).Find(&user).Error

	return &user, err
}

func (s *ServiceUser) CreateUser(user interface{}) error {
	u := user.(*model.User)

	err := storage.DB.Create(u)

	return err.Error
}

func (s *ServiceUser) GetUserByID(userID string) (interface{}, error) {
	var user model.User

	err := storage.DB.Where("user_id = ?", userID).Find(&user).Error

	return &user, err
}

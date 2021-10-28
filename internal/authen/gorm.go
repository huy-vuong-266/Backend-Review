package authen

import (
	"Backend-Review/model"
	"Backend-Review/service"
	storage "Backend-Review/storage"
	"fmt"

	"github.com/jinzhu/gorm"
)

type ServiceAuthen struct{}

func NewServiceAuthen() service.AuthenServiceInterface {
	return &ServiceAuthen{}
}

func (s *ServiceAuthen) CreateToken(token interface{}) error {
	t, ok := token.(*model.Token)
	if !ok {
		return fmt.Errorf("cant parse token")
	}

	err := storage.DB.Where("user_id = ?", t.UserID).Save(t).Error

	return err
}

func (s *ServiceAuthen) CheckIfTokenExist(token string) bool {
	var tokModel model.Token

	var isExist bool = true

	query := storage.DB.Where("token = ?", token)

	err := query.First(&tokModel).Error

	if gorm.IsRecordNotFoundError(err) {
		isExist = false
	} else {
		isExist = true
	}
	return isExist
}

func (s *ServiceAuthen) GetUserIDByAccesstoken(token string) (string, error) {
	var tokModel model.Token
	query := storage.DB.Where("token = ?", token)

	err := query.First(&tokModel).Error
	if gorm.IsRecordNotFoundError(err) || len(tokModel.UserID.String()) == 0 {
		return "", err
	} else {
		return tokModel.UserID.String(), nil
	}

}

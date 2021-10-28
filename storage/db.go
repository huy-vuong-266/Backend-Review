package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func ConnectDB() (err error) {
	DB, err = gorm.Open("mysql", "root:Quang123Huy@@(localhost)/Backend_Review?charset=utf8&parseTime=True&loc=Local")
	if err != nil {

		return gorm.ErrInvalidSQL
	}

	return nil
}

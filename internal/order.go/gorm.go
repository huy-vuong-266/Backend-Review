package order

import (
	"Backend-Review/model"
	"Backend-Review/service"
	"Backend-Review/storage"
	"fmt"
)

type ServiceOrder struct{}

func NewServiceOrder() service.OrderServiceInterface {
	return &ServiceOrder{}
}

func (s *ServiceOrder) CreateOrder(order interface{}) error {
	o, ok := order.(*model.Order)
	if !ok {
		return fmt.Errorf("cant parse order")
	}

	err := storage.DB.Create(o).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceOrder) UpdateOrder(order interface{}) error {
	o, ok := order.(*model.Order)
	if !ok {
		return fmt.Errorf("cant parse order")
	}

	err := storage.DB.Where("order_id = ?", o.OrderID).Save(o).Error
	if err != nil {
		return err
	}
	return nil
}

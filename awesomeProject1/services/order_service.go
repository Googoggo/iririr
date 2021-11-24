package services

import (
	"awesomeProject1/datamodels"
	"awesomeProject1/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(order *datamodels.Order) error
	InsertOrder(order *datamodels.Order) (orderID int64, err error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo()(map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message)(orderID int64, err error)
}

type OrderService struct{
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService{
	return &OrderService{OrderRepository: repository}
}

func (o *OrderService)GetOrderByID(orderID int64)(order *datamodels.Order, err error){
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService)DeleteOrderByID(orderID int64)(isOk bool){
	isOk = o.OrderRepository.Delete(orderID)
	return isOk
}

func (o *OrderService)UpdateOrder(order *datamodels.Order)(err error){
	return o.OrderRepository.Update(order)
}

func (o *OrderService)InsertOrder(order *datamodels.Order)(orderID int64,err error){
	return o.OrderRepository.Insert(order)
}

func (o *OrderService)GetAllOrder()([]*datamodels.Order, error){
	return o.OrderRepository.SelectAll()
}

func (o *OrderService)GetAllOrderInfo()(map[int]map[string]string, error){
	return o.OrderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message)(orderID int64, err error){
	order := &datamodels.Order{
		UserId: message.UserID,
		ProductId: message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}



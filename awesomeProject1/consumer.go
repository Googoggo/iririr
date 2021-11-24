package main

import (
	"awesomeProject1/common"
	"awesomeProject1/rabbitmq"
	"awesomeProject1/repositories"
	"awesomeProject1/services"
	"fmt"
)

func main(){
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	product := repositories.NewProductManager("product",db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order",db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("lpxlpx")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}

package controllers

import (
	"awesomeProject1/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
	OrderService services.IOrderService
}

func (o *OrderController)Get() mvc.View{
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil{
		o.Ctx.Application().Logger().Debug("查询订单失败")
	}
	return mvc.View{
		Name:"order/view.html",
		Data:iris.Map{
			"order":orderArray,
		},
	}
}

func (o *OrderController) GetDelete(){
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	isOk := o.OrderService.DeleteOrderByID(id)
	if isOk {
		o.Ctx.Application().Logger().Debug("成功啊啊")
	} else {
		o.Ctx.Application().Logger().Debug("错误啊啊")
	}
	o.Ctx.Redirect("order/view")
}
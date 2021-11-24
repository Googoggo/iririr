package main

import (
	"awesomeProject1/backend/web/controllers"
	"awesomeProject1/common"
	"awesomeProject1/repositories"
	"awesomeProject1/services"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main(){
	//1.实例
	app := iris.New()
	//2.设置错误，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./backend/web/views",".html").Layout(
		"shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("assets","./backend/web/assets")
	app.OnAnyErrorCode(func(ctx iris.Context){
		ctx.ViewData("message", ctx.Values().GetStringDefault("message","访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil{
		fmt.Println("error")
	}
	//上下文
	ctx ,cancel := context.WithCancel(context.Background())
	defer cancel()
	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))
	orderRepository := repositories.NewOrderManagerRepository("order",db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	//6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

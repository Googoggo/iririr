package main

import (
	"awesomeProject1/common"
	"awesomeProject1/fronted/middleware"
	"awesomeProject1/fronted/web/controllers"
	//"awesomeProject1/rabbitmq"
	"awesomeProject1/repositories"
	"awesomeProject1/services"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板
	app.HandleDir("/public", "/fronted/web/public")
	//访问生成好的html静态文件
	app.HandleDir("/html", "/fronted/web/htmlProductShow")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {

	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepostiory("user",db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	//rabbitmq := rabbitmq.NewRabbitMQSimple("lpxlpx")
	product := repositories.NewProductManager("product", db)
	proService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order",db)
	orderService := services.NewOrderService(order)
	proproduct := app.Party("/product")
	pro := mvc.New(proproduct)
	proproduct.Use(middleware.AuthConProduct)
	pro.Register(proService,orderService,ctx, /*rabbitmq*/)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("127.0.0.1:8082"),
		iris.WithoutBodyConsumptionOnUnmarshal,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}

package main

import (
	"github.com/kataras/iris"
	"log"
	"net/http"
	"xuandan/controller"
	"xuandan/thecrawler"
)

func main() {

	go thecrawler.AllPage()
	go thecrawler.PartPage(1, 10)
	go thecrawler.OverdueGoods()
	app := iris.New()
	app.Get("/goods", controller.Handler)
	app.Get("/goods/searchItemId",controller.HandlerItemId)
	app.Get("/goods/search", controller.HandlerSearch)
	app.Get("/goods/category", controller.HandlerSearchCategory)
	//创建监听
	err := app.Run(iris.Server(&http.Server{Addr: ":9090"}))
	if err != nil {
		log.Println("创建失败：", err)
		return
	}

}

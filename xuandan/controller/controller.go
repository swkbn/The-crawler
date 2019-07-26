package controller

import (
	"encoding/json"
	"github.com/kataras/iris"
	"log"
	"strconv"
	"xuandan/models"
	"xuandan/thecrawler"
)

//分类展示
func HandlerSearchCategory(ctx iris.Context) {
	fqcat := ctx.FormValue("fqcat")
	var count []int
	newFqcat, _ := strconv.Atoi(fqcat)
	if newFqcat > 17 {
		_, err := ctx.WriteString("参数不正确")
		if err != nil {
			log.Println("写入失败：")
		}
		return
	}
	if newFqcat <	1 {
		_, err := ctx.WriteString("参数不正确")
		if err != nil {
			log.Println("写入失败：")
		}
		return
	}
	db := models.Session
	var goods []thecrawler.Goods
	page := ctx.FormValue("page")
	db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("fqcat = ?", newFqcat).Pluck("id", &count)
	act,_:=strconv.Atoi(page)
	count = paging(act, count)
	for _, val := range count {
		var goodsinform models.GoodsItem
		db.Model(&goodsinform).Where("id=?", val).Find(&goodsinform)
		newGoods := thecrawler.FormatConversion(goodsinform)
		goods = append(goods, newGoods)
	}
	newGoods, err := json.Marshal(goods)
	if err != nil {
		log.Println("转换失败")
	}
	_, err = ctx.Write(newGoods)
	if err != nil {
		log.Println("写入失败")
	}
}

//搜索展示
func HandlerSearch(ctx iris.Context) {
	db := models.Session
	var goodsinform []models.GoodsItem
	SearchValue := ctx.FormValue("searchvalue")
	db.Model(&goodsinform).Where("itemtitle LIKE ?", "%"+SearchValue+"%").Order("itemsale2 desc", true).Where("down_type=?", 0).Find(&goodsinform)
	goods, err := json.Marshal(&goodsinform)
	if err != nil {
		log.Println("json转换失败：", err)
	}
	_, err = ctx.Write(goods)
	if err != nil {
		log.Println("写入错误：", err)
	}
}

//展示
func Handler(ctx iris.Context) {
	var count []int
	var goodsinform []thecrawler.Goods
	db := models.Session
	db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Pluck("id", &count)
	page := ctx.FormValue("page")
	act, _ := strconv.Atoi(page)
	if  act<1 {
		ctx.JSON(map[string]string{"mesg":"参数不正确"})
		return
	}
	count = paging(act, count)
	for _, val := range count {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}
	sss, _ := json.Marshal(goodsinform)
	_, err := ctx.Write(sss)
	if err != nil {
		log.Println("写入失败：", err)
	}
}

//分页处理
func paging(page int, count []int) []int {

	newAct := (page - 1) * 100
	end := page * 100
	if newAct > len(count) {
		page = len(count) - 1 - 100
	}
	if end > len(count) {
		end = len(count) - 1
	}
	count = count[newAct:end]
	return count
}

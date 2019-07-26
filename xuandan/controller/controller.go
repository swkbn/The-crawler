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
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	if newFqcat < 1 {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}

	db := models.Session
	var goods []thecrawler.Goods
	page := ctx.FormValue("page")
	db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("fqcat = ?", newFqcat).Pluck("id", &count)
	act, _ := strconv.Atoi(page)
	if act <= 0 {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	count = paging(act, count)
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
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
		log.Println("写入返回值失败")
	}
}

//搜索展示
func HandlerSearch(ctx iris.Context) {
	db := models.Session
	var goodsinform []models.GoodsItem
	SearchValue := ctx.FormValue("searchvalue")
	db.Model(&goodsinform).Where("itemtitle LIKE ?", "%"+SearchValue+"%").Order("itemsale2 desc", true).Where("down_type=?", 0).Find(&goodsinform)
	if goodsinform==nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	goods, err := json.Marshal(&goodsinform)
	if err != nil {
		log.Println("json转换失败：", err)
	}
	_, err = ctx.Write(goods)
	if err != nil {
		log.Println("写入返回值失败：", err)
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
	if act < 1 {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	count = paging(act, count)
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	for _, val := range count {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}
	_, err := ctx.JSON(goodsinform)
	if err != nil {
		log.Println("写入返回值失败：", err)
	}
}

//根据商品id进行搜索
func HandlerItemId(ctx iris.Context) {
	itemId := ctx.FormValue("itemid")
	var goods models.GoodsItem
	db := models.Session
	db.Model(&goods).Where("down_type=?", 0).Where("Itemid= ?", itemId).Find(&goods)
	if goods.Itemid == "" {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	newGoods := thecrawler.FormatConversion(goods)
	ctx.JSON(newGoods)
}

//分页处理
func paging(page int, count []int) []int {
	newAct := (page - 1) * 100
	end := page * 100
	if newAct > len(count) {
		page = len(count) - 1 - 100
		return nil
	}
	if end > len(count) {
		end = len(count) - 1
	}
	count = count[newAct:end]
	return count
}

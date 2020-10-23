package controllers

import (
	"api/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
)

type UserorderController struct {
	beego.Controller
}

func (this *UserorderController) ShowFinaUserorder() {
	//重点get 这里是nil
	//name:=this.GetSession("userName")
	//if name==nil{
	//	beego.Info("用户没登陆")
	//	this.Redirect("/",302)
	//	return
	//}

	countMoney:=0

	//获取页面
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	//每页长度
	pageSize := 10
	start := (pageIndex - 1) * pageSize

	//回到汇总页面,高级查询
	orm := orm.NewOrm()
	var userorderpages []models.Userorder
	var userorders []models.Userorder
	//orm.QueryTable("articletype").All(&userorders)
	qs := orm.QueryTable("userorder")
	//查询总数
	count, _ := qs.Count() //返回数据条目数   加过滤器

	countPage1 := float64(count) / float64(pageSize)
	//总页数
	countPage := math.Ceil(countPage1)

	qs.Limit(pageSize, start).All(&userorders)
	qs.All(&userorderpages)
	//	All(&articles)

	fristPage := false
	endPage := false
	if pageIndex == 1 {
		fristPage = true
	}
	if pageIndex == int(countPage) {
		endPage = true
	}

	for _,v:=range userorderpages{
		countMoney+=v.AmountMoney
	}

	this.Data["countMoney"] = countMoney
	this.Data["userorders"] = userorders
	this.Data["userorderCount"] = count
	this.Data["countPage"] = countPage
	this.Data["pageIndex"] = pageIndex
	this.Data["fristPage"] = fristPage
	this.Data["endPage"] = endPage
	//this.Data["typeAc"] = typeAc
	//this.Data["selectOpt"]=selectOpt
	//
	////articleTypesi
	//this.Data["layoutTile"]="首页"
	//userName:=this.GetSession("userName");
	//this.Data["username"]=userName
	//渲染
	this.TplName = "ReportDay/index.html"

}

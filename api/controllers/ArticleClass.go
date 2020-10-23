package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"api/models"
	"time"
	"strconv"
)

type ArticleClassController struct {
	beego.Controller
}

func(this *ArticleClassController) Index(){
	this.TplName="ArticleClassInfo/index.html"
}
func(this *ArticleClassController)  ArticleClassFrom()  {
	this.TplName = "ArticleClassInfo/form.html"
}

func(this *ArticleClassController) AddArticleClass()  {
	articleClassName:=this.GetString("ArticleClassName")
	remark:=this.GetString("Remark")
	orm:=orm.NewOrm()
	var articelClass models.ArticelClass
	articelClass.Remark=remark
	articelClass.DelFlag=0
	articelClass.ClassName=articleClassName
	articelClass.CreateDate=time.Now()
	articelClass.CreateUserId=this.GetSession("userId").(int)
	articelClass.ParentId=0
	orm.Insert(&articelClass)
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}

func(this *ArticleClassController) GetArticleClassInfo(){
		orm:=orm.NewOrm()
		var articleClass []models.ArticelClass
		orm.QueryTable("articel_class").Filter("parent_id",0).All(&articleClass)
		this.Data["json"]=map[string]interface{}{"rows":articleClass}
		this.ServeJSON()
}

func(this *ArticleClassController)  ShowChildClass()  {
	id:=this.GetString("id")
	beego.Info("id=====>",id)
	idNum,_:=strconv.Atoi(id)
	orm:=orm.NewOrm()
	var articeCla []models.ArticelClass
	orm.QueryTable("articel_class").Filter("parent_id",idNum).All(&articeCla)
	this.Data["json"]=map[string]interface{}{"rows":articeCla}
	this.ServeJSON()
}

//添加子目录
func(this *ArticleClassController)  ShowAddChildClass(){
	cId:=this.GetString("cId")
	id,_:=strconv.Atoi(cId)
	orm:=orm.NewOrm()
	var articelClass models.ArticelClass
	orm.QueryTable("articel_class").Filter("id",id).One(&articelClass)

	this.Data["articelClass"]=articelClass
	this.TplName="ArticleClassInfo/articelForm.html"

	//this.Data["json"]=map[string]interface{}{"flag":"ok"}
	//this.ServeJSON()
	//
	//

}
//AddArticleChildClass
func(this *ArticleClassController) AddArticleChildClass()  {
	articleClassName:=this.GetString("ArticleClassName")
	remark:=this.GetString("Remark")
	id,_:=strconv.Atoi(this.GetString("id"))
	orm:=orm.NewOrm()
	var articelClass models.ArticelClass
	articelClass.Remark=remark
	articelClass.DelFlag=0
	articelClass.ClassName=articleClassName
	articelClass.CreateDate=time.Now()
	articelClass.CreateUserId=this.GetSession("userId").(int)
	articelClass.ParentId=id
	orm.Insert(&articelClass)
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
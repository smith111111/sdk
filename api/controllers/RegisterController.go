package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"crypto/md5"
	"sync"
	"api/models"
	"time"
)

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) ShowRegister() {
	this.TplName = "index.html"

}

func (this *RegisterController) HandlerPost() {
	var mutex sync.Mutex
	mutex.Lock() //加互斥锁
	defer mutex.Unlock() //解互斥锁
	name := this.GetString("mobile")
	pwd := this.GetString("password")
	pwd2 := this.GetString("password2")
	str:=""
	//password2
	//todo后台用户名重复
	if name == "" || pwd == "" {
		beego.Info("用户名或密码为空")
		str="用户名或密码为空"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	if pwd!=pwd2 {
		beego.Info("确认的密码不一致")
		str="确认的密码不一致"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	orm := orm.NewOrm()
	userTwo:=models.UserInfo{}
	orm.QueryTable("user_info").Filter("user_name",name).One(&userTwo)
	if  userTwo.Id !=0{
		beego.Info("注册账号已经存在")
		str="注册账号已经存在"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	//是否存在，插入不同的区
	user := models.UserInfo{}
	user.UserName=name
	data := []byte(pwd)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	user.UserPwd=md5str
	user.AddDate=time.Now()
	user.ModifDate=time.Now()
	fmt.Println("user==============",user)
	_, err := orm.Insert(&user)
	if err != nil {
		beego.Info(err)
		str="注册失败"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	this.Data["json"]=map[string]interface{}{"info":"ok"}
	this.ServeJSON()
}

func (this *RegisterController) HandlerUpdatePwd()  {
	var mutex sync.Mutex
	mutex.Lock() //加互斥锁
	defer mutex.Unlock() //解互斥锁
	name := this.GetString("loginCode")
	pwd := this.GetString("password")
	pwd2 := this.GetString("password2")
	str:=""
	//password2
	//todo后台用户名重复
	if name == "" || pwd == "" {
		beego.Info("密码不能为空")
		str="密码不能为空"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	if pwd!=pwd2 {
		beego.Info("确认的密码不一致")
		str="确认的密码不一致"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	orm := orm.NewOrm()
	userTwo:=models.UserInfo{}
	orm.QueryTable("user_info").Filter("user_name",name).One(&userTwo)
	if  userTwo.Id ==0{
		beego.Info("修改账号不存在")
		str="修改账号不存在"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	//是否存在，插入不同的区
	userTwo.UserName=name
	data := []byte(pwd)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	userTwo.UserPwd=md5str
	userTwo.AddDate=time.Now()
	userTwo.ModifDate=time.Now()
	fmt.Println("userTwo==============",userTwo)
	_, err := orm.Update(&userTwo)
	if err != nil {
		beego.Info(err)
		str="更新失败"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	this.Data["json"]=map[string]interface{}{"info":"ok"}
	this.ServeJSON()
}
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"api/models"
	"crypto/md5"
	"encoding/hex"
	"time"
)

type HomeController struct {
	beego.Controller
}

func(this *HomeController) UserLogin()  {
	//存放id
	userName:=this.GetString("LoginCode")
	LoginPwd:=this.GetString("LoginPwd")
	str:=""
	if userName == "" || LoginPwd == "" {
		beego.Info("用户名或密码为空")
		str="用户名或密码为空"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	h := md5.New()
	h.Write([]byte(LoginPwd)) // 需要加密的字符串为 sharejs.com
	has:= hex.EncodeToString(h.Sum([]byte(nil)))
	beego.Info(has)
	o:=orm.NewOrm()

	var userInfoTwo models.UserInfo
	o.QueryTable("user_info").Filter("user_name",userName).One(&userInfoTwo)
	if userInfoTwo.Id==0{
		beego.Info("登录用户名不存在")
		str="登录用户名不存在"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	var userInfo models.UserInfo
	o.QueryTable("user_info").Filter("user_name",userName).Filter("user_pwd",has).One(&userInfo)
	if userInfo.Id>0{
		this.SetSession("userId",userInfo.Id)
		this.SetSession("userName",userInfo.UserName)
		this.Data["json"]=map[string]interface{}{"flag":"ok"}
	}else{
		beego.Info("密码错误")
		str="密码错误"
		this.Data["json"]=map[string]interface{}{"info":str}
		this.ServeJSON()
		return
	}
	//记住用户名
	check:=this.GetString("remember")
	beego.Info("CH===========",check)
	if check == "on"{
		this.Ctx.SetCookie("userName",userName,time.Second*3600)
		this.Ctx.SetCookie("passWord",LoginPwd,time.Second*3600)
	}else {
		this.Ctx.SetCookie("userName","",-1)
		this.Ctx.SetCookie("passWord","",-1)
	}
	this.ServeJSON()
}
//
func(this *HomeController) Index()  {
	this.TplName="Home/index.html"
	// this.TplName="Home/showIndex.html"
}
func (this *HomeController)IndexShow()  {
	this.TplName="Register/index.html"
}

func (this *HomeController)GetMenus() {
	userId := this.GetSession("userId")
	if userId == nil {
		this.Redirect("/", 302)
		return
	}
	var userInfo models.UserInfo
	orm := orm.NewOrm()
	//根据用户查角色，
	orm.QueryTable("user_info").Filter("id", userId).One(&userInfo)
	orm.LoadRelated(&userInfo, "Roles")

	var actions []*models.ActionInfo
	//遍历角色，获取权限
	for _, roles := range userInfo.Roles {
		orm.LoadRelated(roles, "Actions")
		for _, action := range roles.Actions {
			//获取每个角色的权限
			actions = append(actions, action)
		}
	}
	//判断是否菜单权限
	var roleActions []*models.ActionInfo //1
	for _, act := range actions {
		if act.ActionTypeEnum == 1 {
			roleActions = append(roleActions, act)
		}
	}

	var newActions []models.UserAction
	//获取用户的权限//不禁用的权限   这里建议后面统一过滤，否则后面的过滤不到   .Filter("is_pass",1)
	orm.QueryTable("user_action").Filter("users_id", userId).All(&newActions)
	var newActs []*models.ActionInfo
	if len(newActions) > 0 {
		//过滤菜单权限
		for _, newAction := range newActions {
			var newActi models.ActionInfo
			//根据切片查权限表查询
			orm.QueryTable("action_info").Filter("action_type_enum", 1).Filter("id", newAction.Actions.Id).
				One(&newActi)
			if newActi.Id > 0 {
				newActs = append(newActs, &newActi)
			}
		}
	}
	//合并去重
	newActs = append(newActs, roleActions...)
	newArr := RemoveRepeatedElement(newActs)

	////去除禁用的
	var userForActions []models.UserAction
	orm.QueryTable("user_action").Filter("is_pass", 0).Filter("users_id", userId).All(&userForActions)
	////判断是否找到了登录用户的禁用权限。
	if len(userForActions) > 0 {
		//找到了禁用权限后，清除。
		var newTmep []*models.ActionInfo
		for i, action := range newArr {
			//判断权限的编号是否在禁用的集合中存在，如果返回的是false,表示该权限没有被禁用，存在newTemp切片中。***
			if CheckUserForAction(userForActions, action.Id) == false {
				newTmep = append(newTmep, newArr[i])
			}
		}
	}
		//返回
		this.Data["json"] = map[string]interface{}{"menus": newArr}
		this.ServeJSON()
	}

func RemoveRepeatedElement(arr []*models.ActionInfo) (newArr []*models.ActionInfo) {
	newArr = make([]*models.ActionInfo, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i].Id == arr[j].Id {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}


//判断用户是否有禁用的权限
func CheckUserForAction(userForActions[]models.UserAction,actionId int )(b bool)  {
	b=false
	for i:=0;i<len(userForActions) ;i++  {
		if userForActions[i].Actions.Id==actionId{
			b=true
			break
		}

	}
	return
}
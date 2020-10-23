package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"api/models"
	"strconv"
	"strings"
	"time"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type UserInfoIndex struct {
	beego.Controller
}

type SearchData struct {
	pageIndex  int
	pageSize   int
	txtSname   string
	txtSremark string
	pageConut  int64
}

func (c *UserInfoIndex) UserInfoIndex() {

	c.TplName = "UserInfo/index.html"
}

//addUser
func (this *UserInfoIndex) AddUser() {
	var userInfo = models.UserInfo{}
	o := orm.NewOrm()
	id := this.GetString("id")
	if id != "" {
		idNum1, _ := strconv.Atoi(id)
		userInfo.Id = idNum1
		o.Read(&userInfo)
		userInfo.UserName = this.GetString("userName") //接收用户名
		pwd:= this.GetString("userPwd")
		data := []byte(pwd)
		has := md5.Sum(data)
		md5str := fmt.Sprintf("%x", has)
		userInfo.UserPwd=md5str
		userInfo.Remark = this.GetString("userRemark")
		userInfo.ModifDate = time.Now()
		_, err := o.Update(&userInfo)
		if err == nil {
			//Data中的key必须为"json"
			this.Data["json"] = map[string]interface{}{"flag": "ok", "id": id}

		} else {
			this.Data["json"] = map[string]interface{}{"flag": "no", "id": id}
		}

	} else {
		userInfo.UserName = this.GetString("userName") //接收用户名
		pwd:=this.GetString("userPwd")
		userTwo:=models.UserInfo{}
		o.QueryTable("user_info").Filter("user_name",userInfo.UserName).One(&userTwo)
		if  userTwo.Id !=0{
			beego.Info("新增账号已经存在")
			str:="新增账号已经存在"
			this.Data["json"]=map[string]interface{}{"flag": str, "id": id}
			this.ServeJSON()
			return
		}
      //加密MD5
      //为什么要对密码做两次MD5
      //现在存在的一些反查md5的软件，做两次为了更好的保密
		 h := md5.New()
		 h.Write([]byte(pwd)) // 需要加密的字符串为 sharejs.com
		has:= hex.EncodeToString(h.Sum([]byte(nil)))
		beego.Info("ha",has)
		userInfo.UserPwd=has
		userInfo.Remark = this.GetString("userRemark")
		userInfo.ModifDate = time.Now()
		userInfo.AddDate = time.Now()
		userInfo.DelFlag = 0 //表示正常，1表示表示软删除。
		_, err := o.Insert(&userInfo)
		if err == nil {
			//Data中的key必须为"json"
			this.Data["json"] = map[string]interface{}{"flag": "ok", "id": id}

		} else {
			this.Data["json"] = map[string]interface{}{"flag": "no", "id": id}
		}
	}
	//怎样将数据生成JSON.
	this.ServeJSON()
}

func (this *UserInfoIndex) UserInfoIndexData() {
	pageIndex, _ := this.GetInt("page")
	pageSize, _ := this.GetInt("rows")
	txtSname := this.GetString("txtSname")
	txtSremark := this.GetString("txtSremark")

	var searchData = SearchData{}
	searchData.pageIndex = pageIndex
	searchData.pageSize = pageSize
	searchData.txtSname = txtSname
	searchData.txtSremark = txtSremark
	userSear := searchData.GetUserData()

	this.Data["json"] = map[string]interface{}{"rows": userSear, "total": searchData.pageConut}
	this.ServeJSON()
}

func (this *SearchData) GetUserData() []models.UserInfo {
	orm := orm.NewOrm()
	start := (this.pageIndex - 1) * this.pageSize
	var users []models.UserInfo
	qs := orm.QueryTable("user_info")
	//icontains模糊查询该写入的字段
	if this.txtSname != "" {
		qs = qs.Filter("user_name__icontains", this.txtSname)
	}
	if this.txtSremark != "" {
		qs = qs.Filter("remark__icontains", this.txtSremark)
	}
	qs = qs.Filter("del_flag", 0)
	this.pageConut, _ = qs.Count()
	//降序用 -Id
	qs.OrderBy("Id").Limit(this.pageSize, start).All(&users)
	//

	return users
}

func (this *UserInfoIndex) DeleteUser() {
	strArray := this.GetString("ostrArr")

	strNum1 := strings.Split(strArray, ",")
	orm := orm.NewOrm()
	var userInfo = models.UserInfo{}
	for i := 0; i < len(strNum1); i++ {
		iNum2, _ := strconv.Atoi(strNum1[i])
		userInfo.Id = iNum2
		orm.Delete(&userInfo)
	}

	this.Data["json"] = map[string]interface{}{
		"flag": "ok"}

	this.ServeJSON()

}

func (this *UserInfoIndex) UpdateUser() {

	//
	//if err == nil {
	//	//Data中的key必须为"json"
	//	this.Data["json"] = map[string]interface{}{"flag": "ok"}
	//
	//} else {
	//	this.Data["json"] = map[string]interface{}{"flag": "no"}
	//}
	//this.ServeJSON()
}

func (this *UserInfoIndex) FindUserById() {
	id, _ := this.GetInt("Id")
	orm := orm.NewOrm()
	userInfo := models.UserInfo{}
	userInfo.Id = id
	err := orm.Read(&userInfo)
	if err == nil {
		this.Data["json"] = map[string]interface{}{
			"userName":   userInfo.UserName,
			"userPwd":    userInfo.UserPwd,
			"userRemark": userInfo.Remark,
			"id":         userInfo.Id,
		}
	} else {
		this.Data["json"] = map[string]interface{}{"UserInfo": ""}
	}
	this.ServeJSON()
}

//查询角色分配页面
func (this *UserInfoIndex) FindUserRoleById() {
	id, _ := this.GetInt("id")
	var userInfo = models.UserInfo{}
	var roleInfos []models.RoleInfo
	var roleNowInfos []*models.RoleInfo
	userInfo.Id = id
	orm := orm.NewOrm()
	//查询当前角色
	orm.Read(&userInfo)

	//查当前用户有的角色
	orm.LoadRelated(&userInfo, "Roles")

	for _, role := range userInfo.Roles {
		roleNowInfos = append(roleNowInfos, role)
	}

	//查询所有角色
	orm.QueryTable("role_info").Filter("del_flag", 0).OrderBy("Id").All(&roleInfos)

	this.Data["userInfo"] = userInfo
	this.Data["roleInfos"] = roleInfos
	this.Data["roleNowInfos"] = roleNowInfos

	this.TplName = "UserInfo/UserRoleMenu.html"
}

func (this *UserInfoIndex) SetUserRole() {
	userId, _ := this.GetInt("userId")
	var roleIdList []int
	allkeys := this.Ctx.Request.PostForm
	for key, _ := range allkeys {
		if strings.Contains(key, "cba_") {
			id := strings.Replace(key, "cba_", "", -1)
			roleId, _ := strconv.Atoi(id)
			roleIdList = append(roleIdList, roleId)
		}
	}
	//查询用户信息
	o := orm.NewOrm()
	var userInfo models.UserInfo
	o.QueryTable("user_info").Filter("Id", userId).One(&userInfo)
	//查询用户具有的角色信息
	o.LoadRelated(&userInfo, "Roles")
	m2m := o.QueryM2M(&userInfo, "Roles")
	//删除用户已经有的角色信息
	o.Begin() //开启事务
	var err1 error
	var err2 error
	for _, role := range userInfo.Roles {
		_, err1 = m2m.Remove(role)
	}
	//重新给用户分配角色信息
	var roleInfo models.RoleInfo
	for i := 0; i < len(roleIdList); i++ {
		o.QueryTable("role_info").Filter("Id", roleIdList[i]).One(&roleInfo)
		_, err2 = m2m.Add(roleInfo)
	}
	if err2 != nil || err1 != nil {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"flag": "no"}
	} else {

		o.Commit()
		this.Data["json"] = map[string]interface{}{"flag": "yes"}
	}
	this.ServeJSON()

}

//为用户分配权限
func (this *UserInfoIndex)ShowUserAction()  {
	//接收用户Id
	userId,_:=strconv.Atoi(this.GetString("userId"))
	o:=orm.NewOrm()
	var userInfo models.UserInfo
	//查询用户信息
	o.QueryTable("user_info").Filter("Id",userId).One(&userInfo)
	//查询出用户已经有的权限编号。
	var userExtActions []models.UserAction//一对多直接查
	o.QueryTable("user_action").Filter("users_id",userId).All(&userExtActions)
	//查询出所有的权限信息
	var allActionList []models.ActionInfo
	o.QueryTable("action_info").Filter("del_flag",0).All(&allActionList)
	beego.Info(userInfo)
	beego.Info(allActionList)
	beego.Info(userExtActions)
	this.Data["userInfo"]=userInfo
	this.Data["allActions"]=allActionList
	this.Data["userExtActions"]=userExtActions
	this.TplName="UserInfo/ShowSetUserAction.html"
}

//完成用户权限的分配
func (this *UserInfoIndex)SetUserAction()  {
	actionId,_:=strconv.Atoi(this.GetString("actionId"))
	isPass,_:=this.GetBool("isPass")
	userId,_:=strconv.Atoi(this.GetString("userId"))
	var isExt int
	if isPass{
		isExt=1
	}else{
		isExt=0
	}
	o:=orm.NewOrm()
	var userAction models.UserAction
	o.QueryTable("user_action").Filter("users_id",userId).Filter("actions_id",actionId).One(&userAction)
	if userAction.Id>0{
		//如果用户有权限，直接修改
		userAction.IsPass=isExt
		o.Update(&userAction)
	}else{
		//如果没有权限，直接添加
		var actionInfo models.ActionInfo
		o.QueryTable("action_info").Filter("Id",actionId).One(&actionInfo)
		var userInfo models.UserInfo
		o.QueryTable("user_info").Filter("Id",userId).One(&userInfo)
		userAction.IsPass=isExt
		userAction.Actions=&actionInfo
		userAction.Users=&userInfo
		o.Insert(&userAction)

	}
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}

//删除用户权限
func(this *UserInfoIndex) DeleteUserAction()  {
	actionId,_:=strconv.Atoi(this.GetString("actionId"))
	userId,_:=strconv.Atoi(this.GetString("userId"))
	o:=orm.NewOrm()
	var userAciton models.UserAction
	o.QueryTable("user_action").Filter("users_id",userId).Filter("actions_id",actionId).One(&userAciton)
	o.Delete(&userAciton)

	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}



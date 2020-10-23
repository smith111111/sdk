package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"api/models"
	"strconv"
	"strings"
	"time"
)

type RoleInfoIndex struct {
	beego.Controller
}

type RoleSearch struct {
	pageIndex      int
	pageSize       int
	roleName       string
	roleRemark     string
	roleTotalCount int64
}

func (this *RoleInfoIndex) RoleInfoIndex() {
	beego.Info("进入")
	this.TplName = "RoleInfo/index.html"
}

func (this *RoleInfoIndex) RoleInfoIndexData() {
	pageIndex, _ := this.GetInt("page") //当前一
	pageSize, _ := this.GetInt("rows")  //每页数

	roleName := this.GetString("roleName")     //角色名称
	roleRemark := this.GetString("roleRemark") //角色备注
	var rolesh = RoleSearch{}
	rolesh.pageIndex = pageIndex
	rolesh.pageSize = pageSize
	rolesh.roleName = roleName
	rolesh.roleRemark = roleRemark

	roles := rolesh.SearchData()
	this.Data["json"] = map[string]interface{}{
		"rows": roles, "total": rolesh.roleTotalCount,
	}
	this.ServeJSON()
}

func (this *RoleSearch) SearchData() []models.RoleInfo {
	start := (this.pageIndex - 1) * this.pageSize
	var roles []models.RoleInfo
	orm := orm.NewOrm()

	temp := orm.QueryTable("role_info")

	if this.roleName != "" {
		temp = temp.Filter("role_name__icontains", this.roleName)
	}
	if this.roleRemark != "" {
		temp = temp.Filter("remark__icontains", this.roleRemark)
	}
	temp = temp.Filter("del_flag", 0)

	this.roleTotalCount, _ = temp.Count()
	temp.Limit(this.pageSize, start).OrderBy("Id").All(&roles)
	return roles
}

func (this *RoleInfoIndex) RoleInfoFormUrl() {
	id, _ := this.GetInt("id")
	roleName := this.GetString("roleName")
	roleRemark := this.GetString("remark")
	this.Data["roleName"] = roleName
	this.Data["roleRemark"] = roleRemark
	this.Data["id"] = id
	this.TplName = "RoleInfo/updateForm.html"
}

func (this *RoleInfoIndex) RoleAddInfoFormUrl() {
	this.TplName = "RoleInfo/form.html"
}
func (this *RoleInfoIndex) RoleInfoAdd() {
	roleName := this.GetString("RoleName")
	roleRemark := this.GetString("Remark")

	orm := orm.NewOrm()
	var role models.RoleInfo
	role.RoleName = roleName
	role.Remark = roleRemark
	role.AddDate = time.Now()
	role.DelFlag = 0
	role.ModifDate = time.Now()

	orm.Insert(&role)
	this.Data["json"] = map[string]interface{}{
		"flag": "ok",
	}
	this.ServeJSON()

}

func (this *RoleInfoIndex) DeleteRole() {
	param := this.GetString("param")
	params := strings.Split(param, ",")

	orm := orm.NewOrm()
	var role = models.RoleInfo{}
	for i := 0; i < len(params); i++ {
		role.Id, _ = strconv.Atoi(params[i])
		orm.Delete(&role)
	}

	this.Data["json"] = map[string]interface{}{
		"flag": "ok",
	}
	this.ServeJSON()

}

func (this *RoleInfoIndex) SearRole() {
	id, _ := this.GetInt("Id")
	orm := orm.NewOrm()
	var role = models.RoleInfo{}
	role.Id = id
	orm.Read(&role)
	beego.Info(role)
	this.Data["json"] = map[string]interface{}{
		"role": role,
	}
	this.ServeJSON()
}

func (this *RoleInfoIndex) UpdateRole() {
	id, _ := this.GetInt("roleId")
	roleName := this.GetString("roleName")
	roleRemark := this.GetString("roleRemark")
	beego.Info(id)
	beego.Info(roleName)
	beego.Info(roleRemark)
	var role models.RoleInfo
	role.Id = id
	orm := orm.NewOrm()
	orm.Read(&role)
	role.Remark = roleRemark
	role.RoleName = roleName
	_, err := orm.Update(&role)
	if err == nil {
		this.Data["json"] = map[string]interface{}{
			"flag": "ok",
		}
	}
	this.ServeJSON()
}

//角色分配权限

func (this *RoleInfoIndex)ShowSetRoleAction()  {
	//1:获取角色编号
	roleId,_:=strconv.Atoi(this.GetString("roleId"))
	//2:根据该角色编号查询出对应的角色信息
	o:=orm.NewOrm()
	var roleInfo models.RoleInfo
	o.QueryTable("role_info").Filter("Id",roleId).One(&roleInfo)
	//3:查询该角色已经有的权限信息
	var actions []*models.ActionInfo
	o.LoadRelated(&roleInfo,"Actions")
	for _,action:=range roleInfo.Actions {
		actions=append(actions,action)
	}
	//4:查询出所有的权限信息
	var allActionList []models.ActionInfo
	o.QueryTable("action_info").Filter("del_flag",0).All(&allActionList)
	//5:展示出对应的数据
	this.Data["roleInfo"]=roleInfo
	this.Data["roleExtActions"]=actions
	this.Data["allActionList"]=allActionList
	this.TplName="RoleInfo/ShowSetRoleAction.html"
}

func (this *RoleInfoIndex) SetRoleAction()  {
	//获取角色编号
	roleId,_:=strconv.Atoi(this.GetString("roleId"))
	allKeys:=this.Ctx.Request.PostForm
	var actionIdList[]int
	for key,_:= range allKeys{
		if strings.Contains(key,"cba_") {
			id:=strings.Replace(key,"cba_","",-1)
			actionId,_:=strconv.Atoi(id)
			actionIdList=append(actionIdList,actionId)
		}
	}
	//查询角色的信息
	var roleInfo models.RoleInfo
	o:=orm.NewOrm()
	o.QueryTable("role_info").Filter("Id",roleId).One(&roleInfo)
	//查询出角色对应的权限
	o.LoadRelated(&roleInfo,"Actions")
	m2m:=o.QueryM2M(&roleInfo,"Actions")
	o.Begin()
	var err1 error
	var err2 error
	for _,action:=range roleInfo.Actions {
		_,err1=m2m.Remove(action)
	}
	var actionInfo models.ActionInfo
	for i:=0;i<len(actionIdList) ;i++  {
		o.QueryTable("action_info").Filter("Id",actionIdList[i]).One(&actionInfo)
		_,err2=m2m.Add(actionInfo)
	}
	if err1!=nil||err2!=nil{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
		o.Rollback()
	}else{
		this.Data["json"]=map[string]interface{}{"flag":"yes"}
		o.Commit()
	}
	this.ServeJSON()
}
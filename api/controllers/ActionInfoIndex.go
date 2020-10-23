package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"github.com/astaxie/beego/orm"
	"api/models"
	"path"
	"time"
	"os"
	"strings"
)

type ActionInfoIndex struct {
	beego.Controller
}

func(this *ActionInfoIndex) ActionInfoIndex(){
	this.TplName = "ActionInfo/index.html"
}

func(this *ActionInfoIndex) GetActionInfo()  {
	pageIndex,_:=strconv.Atoi(this.GetString("page"))
	pageSize,_:=strconv.Atoi(this.GetString("rows"))
	start:=(pageIndex-1)*pageSize
	o:=orm.NewOrm()
	var actions []models.ActionInfo
	o.QueryTable("action_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&actions)
	count,_:=o.QueryTable("action_info").Filter("del_flag",0).Count()
	this.Data["json"]=map[string]interface{}{"rows":actions,"total":count}
	this.ServeJSON()
}

func (this *ActionInfoIndex)  ActionAdd(){
	this.TplName="ActionInfo/addActionInfo.html"
}

func (this *ActionInfoIndex)ActionFileUp()  {
	f,h,err:=this.GetFile("fileUp")
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"flag":"no","msg":"上传文件错误"}
	}else{
		//获取上传的文件名称
		fileName:=h.Filename
		//获取上传文件的类型
		fileExt:=path.Ext(fileName)
		if fileExt==".jpg"||fileExt==".png"{
			if h.Size<50000000 {//判断上传文件的大小
				//创建上传图片文件存放的路径。
				dirPath:="./static/fileUp/"+strconv.Itoa(time.Now().Year())+"/"+time.Now().Month().String()+"/"+strconv.Itoa(time.Now().Day())+"/"
				_,err:=os.Stat(dirPath)
				if err!=nil{//表示没有目录信息
					dirError:=os.MkdirAll(dirPath,os.ModePerm)//创建目录
					if dirError!=nil{
						this.Data["json"]=map[string]interface{}{"flag":"no","msg":"目录创建失败!!"}
						return
					}
				}
				//按照日期时间对文件进行重命名。
				fileNewName:=strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+strconv.Itoa(time.Now().Hour())+strconv.Itoa(time.Now().Minute())+strconv.Itoa(time.Now().Nanosecond())//获取毫秒数
				fullDir:=dirPath+fileNewName+fileExt//构建完整的路径
				fileErr:=this.SaveToFile("fileUp",fullDir)//进行文件的保存
				if fileErr==nil{
					this.Data["json"]=map[string]interface{}{"flag":"ok","msg":fullDir}


				}else{
					this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件上传失败！！"}
				}

			}else{
				this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件太大"}
			}

		}else{
			this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件类型错误!!"}
		}
	}
	defer f.Close()
	this.ServeJSON()
}

func(this *ActionInfoIndex) AddAction()  {
	var actionInfo=models.ActionInfo{}
	actionInfo.DelFlag=0
	actionInfo.ModifDate=time.Now()
	actionInfo.AddDate=time.Now()
	actionInfo.Remark=this.GetString("Remark")
	actionInfo.MenuIcon=this.GetString("MenuIcon")
	actionInfo.Url=this.GetString("Url")
	actionInfo.ActionInfoName=this.GetString("ActionInfoName")
	actionInfo.ActionTypeEnum,_=strconv.Atoi(this.GetString("ActionTypeEnum"))
	actionInfo.IconWidth=0
	actionInfo.IconHeight=0
	actionInfo.HttpMethod=this.GetString("HttpMethod")
	o:=orm.NewOrm()
	_,err:=o.Insert(&actionInfo)
	if err==nil {
		this.Data["json"]=map[string]interface{}{"flag":"ok"}
	}else{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()
}

func (this *ActionInfoIndex) DeleteAction() {
	strArray := this.GetString("ostrArr")

	strNum1 := strings.Split(strArray, ",")
	orm := orm.NewOrm()
	var actionInfo = models.ActionInfo{}
	for i := 0; i < len(strNum1); i++ {
		iNum2, _ := strconv.Atoi(strNum1[i])
		actionInfo.Id = iNum2
		orm.Delete(&actionInfo)
	}

	this.Data["json"] = map[string]interface{}{
		"flag": "ok"}

	this.ServeJSON()

}
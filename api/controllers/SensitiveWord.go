package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"github.com/astaxie/beego/orm"
	"api/models"
)

type SensitiveWordController struct {
	beego.Controller
}

func(this *SensitiveWordController) Index()  {
   this.TplName="SensitiveWord/index.html"
}
//添加敏感词
func (this *SensitiveWordController)AddWords()  {
	contentMsg:=this.GetString("contentMsg")
	strs:=strings.Split(contentMsg,"\r\n")
	var strMsg[]models.SensitiveWord
	for i:=0;i<len(strs) ;i++  {
		var sensitiveWord=models.SensitiveWord{}
		words:=strings.Split(strs[i],"=")
		sensitiveWord.WordPattern=words[0]
		if words[1]=="{BANNED}"{
			sensitiveWord.IsForbid=1
		}else  if words[1]=="{MOD}"{
			sensitiveWord.IsMod=1
		}else{
			sensitiveWord.ReplaceWord=words[1]
		}
		strMsg=append(strMsg,sensitiveWord)

	}
	o:=orm.NewOrm()
	o.InsertMulti(len(strMsg),strMsg)
	this.Data["json"]=map[string]interface{}{"flag":"yes"}
	this.ServeJSON()

}

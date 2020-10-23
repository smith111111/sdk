package controllers

import (
	"github.com/astaxie/beego"

	"github.com/dchest/captcha"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) GetInfo() {
	name := c.Ctx.GetCookie("userName")
	if name != ""{
		c.Data["userName"] = name
		c.Data["check"] = "checked"//最好动态获取
	}
	c.Data["data"] = "aaa"

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "scloudrun@gmail.com"
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	c.Data["CaptchaId"] = d.CaptchaId
	c.TplName = "index.html"
}

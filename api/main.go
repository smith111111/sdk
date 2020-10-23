package main

import (
	"api/models"
	_ "api/models"
	_ "api/routers"
	"api/utils"
	"github.com/astaxie/beego"
)

func CheckId(roles []*models.RoleInfo, roleId int) (flag bool) {
	flag = false
	for i := 0; i < len(roles); i++ {
		if roles[i].Id == roleId {
			flag = true
			break
		}
	}
	return
}
func  CheckAction(userExtActionList[]*models.ActionInfo,actionId int)(b bool)  {
	b=false
	for i:=0;i<len(userExtActionList) ;i++  {
		if userExtActionList[i].Id==actionId {
			b=true
			break
		}
	}
	return
}

func CheckUserAction(userExtActionList[]models.UserAction,actionId int)(b bool) {
	b=false
	for i:=0;i<len(userExtActionList);i++ {
		if userExtActionList[i].Actions.Id==actionId {
			b=true
			break
		}
	}
	return
}
//判断用户具有某个权限是禁止还是允许
func CheckUserActionId(userExtActionList[]models.UserAction)(b bool)  {
	b=false
	for i:=0;i<len(userExtActionList) ;i++  {
			if userExtActionList[i].IsPass==1 {
				b=true
				break//注意break的位置
			}
	}
	return
}

//初始化旧的仙侠四个区
func initOldXianxia(){
	go func() {
		utils.InitXianXiaOneDataDetailsTask()
	}()

	//go func() {
	//	utils.InitXianXiaTwoDataDetailsTask()
	//}()
	//go func() {
	//	utils.InitXianXiaThreeDataDetailsTask()
	//}()
	//go func() {
	//	utils.InitXianXiaFourDataDetailsTask()
	//}()
   //订单旧仙侠
	go func() {
		utils.InitXianXiaOneOrderDataDetailsTask()
	}()

}

//初始化新的仙侠二个区
func initNewXianxia(){
	go func() {
		utils.InitXianXiaNewOneDataDetailsTask()
	}()

	//go func() {
	//	utils.InitXianXiaNewTwoDataDetailsTask()
	//}()
	//

	//order
	go func() {
		utils.InitXianXiaNewOneOrderDataDetailsTask()
	}()
}

func initGuozhan()  {
	go func() {
		utils.InitGuozhanOneDataDetailsTask()
	}()
	go func() {
		utils.InitGuozhanOneOrderDataDetailsTask()
	}()
}

func initChuanqi()  {
	go func() {
		utils.InitChuanqiOneDataDetailsTask()
	}()
	go func() {
		utils.InitChuanqiOneOrderDataDetailsTask()
	}()
}

func main() {
	initChuanqi()
	initGuozhan()
	//////初始化旧的仙侠四个区
	//initOldXianxia()
	initNewXianxia()
	go func() {
		utils.ReportUserDataCalculation()
	}()
	beego.AddFuncMap("checkId", CheckId)
	beego.AddFuncMap("checkAction",CheckAction)
	beego.AddFuncMap("checkUserAction", CheckUserAction)
	beego.AddFuncMap("checkUserActionId",CheckUserActionId)
	beego.SetStaticPath("/", "views")
	beego.Run()
}



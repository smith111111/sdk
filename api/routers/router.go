package routers

import (

	"api/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/Admin/UserInfo/index", &controllers.UserInfoIndex{}, "get:UserInfoIndex")
	beego.Router("/Admin/UserInfo/AddUser", &controllers.UserInfoIndex{}, "post:AddUser")

	beego.Router("/Admin/UserInfo/GetUserInfo", &controllers.UserInfoIndex{}, "post:UserInfoIndexData")

	beego.Router("/Admin/UserInfo/DeleteUser", &controllers.UserInfoIndex{}, "post:DeleteUser")


	beego.Router("/Admin/UserInfo/UpdateUser", &controllers.UserInfoIndex{}, "post:UpdateUser")
 	beego.Router("/Admin/UserInfo/FindUserById", &controllers.UserInfoIndex{}, "post:FindUserById")
	beego.Router("/Admin/UserInfo/GetUserRoleFormUrl", &controllers.UserInfoIndex{}, "get:FindUserRoleById")
	beego.Router("/Admin/UserInfo/SetUserRole", &controllers.UserInfoIndex{}, "post:SetUserRole")
//展示用户权限
	beego.Router("/Admin/UserInfo/ShowUserAction", &controllers.UserInfoIndex{}, "get:ShowUserAction")
///
	beego.Router("Admin/UserInfo/SetUserAction", &controllers.UserInfoIndex{}, "post:SetUserAction")
	beego.Router("Admin/UserInfo/DeleteUserAction", &controllers.UserInfoIndex{}, "post:DeleteUserAction")

	///Admin/UserInfo/DeleteUserAction
	//-----------------------------Role---------------------------------------------------------
	beego.Router("/Admin/RoleInfo/index", &controllers.RoleInfoIndex{}, "get:RoleInfoIndex")
	beego.Router("/Admin/RoleInfo/GetRoleInfo", &controllers.RoleInfoIndex{}, "post:RoleInfoIndexData")
	//获取角色的url路径
	beego.Router("/Admin/RoleInfo/GetRoleFormUrl", &controllers.RoleInfoIndex{}, "get:RoleInfoFormUrl")
	//
	beego.Router("/Admin/RoleInfo/GetRoleAddFormUrl", &controllers.RoleInfoIndex{}, "get:RoleAddInfoFormUrl")
	beego.Router("/Admin/RoleInfo/AddRole", &controllers.RoleInfoIndex{}, "post:RoleInfoAdd")
	beego.Router("/Admin/RoleInfo/DeleteRole", &controllers.RoleInfoIndex{}, "post:DeleteRole")
	beego.Router("/Admin/RoleInfo/SearRole", &controllers.RoleInfoIndex{}, "post:SearRole")
	beego.Router("/Admin/RoleInfo/UpdateRole", &controllers.RoleInfoIndex{}, "post:UpdateRole")
	//角色分配权限
	beego.Router("/Admin/RoleInfo/ShowSetRoleAction", &controllers.RoleInfoIndex{}, "get:ShowSetRoleAction")
	beego.Router("/Admin/RoleInfo/SetRoleAction", &controllers.RoleInfoIndex{}, "post:SetRoleAction")


	//-------------------权限
	beego.Router("/Admin/ActionInfo/index", &controllers.ActionInfoIndex{}, "get:ActionInfoIndex")
	beego.Router("/Admin/ActionInfo/GetActionInfo",&controllers.ActionInfoIndex{},"post:GetActionInfo")
	beego.Router("/Admin/ActionInfo/GetActionInfoAddFormUrl", &controllers.ActionInfoIndex{}, "get:ActionAdd")
	beego.Router("/Admin/ActionInfo/FileUp",&controllers.ActionInfoIndex{},"post:ActionFileUp")
	beego.Router("/Admin/ActionInfo/AddAction",&controllers.ActionInfoIndex{},"post:AddAction")
	beego.Router("/Admin/ActionInfo/DeleteAction", &controllers.ActionInfoIndex{}, "post:DeleteAction")

//首页
	beego.Router("/Admin/ActionInfo/AddAc",&controllers.HomeController{},"post:GetMenus")
	beego.Router("/Admin/ActionInfo/AddActionShow",&controllers.HomeController{},"post:IndexShow")


//登录

	beego.Router("/Login/UserLogin",&controllers.HomeController{},"post:UserLogin")
	beego.Router("/Admin/Home/Index",&controllers.HomeController{},"get:Index")
	beego.InsertFilter("/Admin/*",beego.BeforeExec,FilterUserAction)


	//--------------------文章类别  ------------
	beego.Router("/Admin/ArticleClass/Index",&controllers.ArticleClassController{},"get:Index")
	//
	beego.Router("/Admin/AticleClassInfo/GetAticleClassUrl",&controllers.ArticleClassController{},"get:ArticleClassFrom")


	beego.Router("/Admin/ArticleClassInfo/AddArticleClass",&controllers.ArticleClassController{},"post:AddArticleClass")

//加载  /Admin/ArticleClassInfo/GetArticleClassInfo
	beego.Router("/Admin/ArticleClassInfo/GetArticleClassInfo",&controllers.ArticleClassController{},"post:GetArticleClassInfo")

	beego.Router("/Admin/ArticelClass/ShowChildClass",&controllers.ArticleClassController{},"post:ShowChildClass")

	beego.Router("/Admin/ArticelClass/ShowAddChildClass",&controllers.ArticleClassController{},"get:ShowAddChildClass")

	beego.Router("/Admin/ArticleClassInfo/AddArticleChildClass",&controllers.ArticleClassController{},"post:AddArticleChildClass")


//----------------------------------新闻
	beego.Router("/Admin/ArticleInfo/GetArticleInfoIndex",&controllers.ArticleInfoController{},"get:GetArticleInfoIndex")
	beego.Router("/Admin/ArticelInfo/GetArticelInfo",&controllers.ArticleInfoController{},"post:GetArticelInfo")
	beego.Router("/Admin/ArticelInfo/ShowAddArticelInfo",&controllers.ArticleInfoController{},"get:ShowAddArticelInfo")

///Admin/ArticeInfo/FileUp
	beego.Router("/Admin/ArticeInfo/FileUp",&controllers.ArticleInfoController{},"post:FileUp")
	beego.Router("/Admin/ArticeInfo/AddArtice",&controllers.ArticleInfoController{},"post:AddArtice")

	beego.Router("/Admin/ArticelInfo/AddComment",&controllers.ArticleInfoController{},"post:AddComment")
	beego.Router("/Admin/ArticelInfo/LoadCommentMsg",&controllers.ArticleInfoController{},"post:LoadCommentMsg")


	//----评论
	//beego.Router("/Admin/SensitiveWord/AddWords",&controllers.SensitiveWordController{},"post:AddWords")
	//beego.Router("/Admin/SensitiveWord/Index",&controllers.SensitiveWordController{},"get:Index")
	//注册
	beego.Router("/Game/Register", &controllers.RegisterController{}, "POST:HandlerPost")
	//修改密码
	beego.Router("/Game/UpdatePwd", &controllers.RegisterController{}, "POST:HandlerUpdatePwd")



    //财务报表
	beego.Router("/Admin/ArticleInfo/GetFinanStateIndex", &controllers.ArticleInfoController{}, "get:GetFinaArticleInfoIndex")
	beego.Router("/Admin/ArticelInfo/GetFinalArticelInfo",&controllers.ArticleInfoController{},"POST:GetFinalArticelInfo")

    //订单报表 /Admin/UserInfo/GetUserOrderInfoIndex

	beego.Router("/Admin/ArticleInfo/GetUserOrderInfoIndex", &controllers.ArticleInfoController{}, "get:GetUserOrderInfoIndex")
	beego.Router("/Admin/ArticleInfo/GetUserOrderInfo",&controllers.ArticleInfoController{},"POST:GetUserOrderInfo")

    //渠道报表
	beego.Router("/Admin/ArticleInfo/GetChannelReportIndex",&controllers.ArticleInfoController{},"get:GetChannelReportIndex")
	beego.Router("/Admin/ArticleInfo/GetUserOrderChanlInfo",&controllers.ArticleInfoController{},"POST:GetUserOrderChanlInfo")


	//留存报表
	beego.Router("/Admin/ArticleInfo/GetRetentionReportIndex",&controllers.ArticleInfoController{},"get:GetRetentionReportIndex")
	beego.Router("/Admin/ArticleInfo/GetRetentionInfo",&controllers.ArticleInfoController{},"POST:GetRetentionInfo")

}
//过滤器   后面可细节到登录 和 按钮显示权限  暂时不需要
//
//func FilterUserAction(ctx *context.Context)  {
//	//userId:=ctx.Input.Session("userId")
//	userName:=ctx.Input.Session("userName")
//	if userName!=nil{
//		path:=ctx.Request.URL.Path
//		method:=ctx.Request.Method
//		//
//		var actionInfo models.ActionInfo
//		orm:=orm.NewOrm()
//		orm.QueryTable("action_info").Filter("url",path).Filter("http_method",method).One(&actionInfo)
//		fmt.Println("123===",userName)
//		if userName=="123"{
//			return
//		}
//		if actionInfo.Id>0{
//			//判断权限
//			//1.查询用户
//			var userInfo models.UserInfo
//			orm.QueryTable("user_info").Filter("user_name",userName).One(&userInfo)
//			//2.查询是否user_action表中有权限关联
//			var userAction models.UserAction
//			orm.QueryTable("user_action").
//				Filter("users_id",userInfo.Id).Filter("actions_id",actionInfo.Id).One(&userAction)
//			if userAction.Id>0{
//				if userAction.IsPass==1{
//					return
//				}else{
//					ctx.Redirect(302,"/")
//				}
//			}else{//不在用户角色的情况下判断，所以下面不判断ispass ****
//				//用户--角色--权限
//				orm.LoadRelated(&userInfo,"roles")
//				var actionFinal []*models.ActionInfo
//				for _,role:=range userInfo.Roles{
//					  orm.LoadRelated(role,"actions")
//					  for _,action:=range role.Actions {
//						  if action.Id == actionInfo.Id{
//						  actionFinal = append(actionFinal, action)
//						  }
//					  }
//				}
//				if len(actionFinal)>1{
//						return
//				}else{
//					ctx.Redirect(302,"/")
//				}
//			}
//		}else{
//			//没有对应的方法
//			ctx.Redirect(302,"/")
//		}
//	}else{
//		//用户没登录
//		ctx.Redirect(302,"/")
//	}
//}
//
//


func FilterUserAction(ctx *context.Context)  {
	//userId:=ctx.Input.Session("userId")
	userName:=ctx.Input.Session("userName")
	if userName!=nil{
		//path:=ctx.Request.URL.Path
		//method:=ctx.Request.Method
		////
		//var actionInfo models.ActionInfo
		//orm:=orm.NewOrm()
		//orm.QueryTable("action_info").Filter("url",path).Filter("http_method",method).One(&actionInfo)
		//fmt.Println("123===",userName)
		//if userName=="123"{
		//	return
		//}
		//if actionInfo.Id>0{
		//	//判断权限
		//	//1.查询用户
		//	var userInfo models.UserInfo
		//	orm.QueryTable("user_info").Filter("user_name",userName).One(&userInfo)
		//	//2.查询是否user_action表中有权限关联
		//	var userAction models.UserAction
		//	orm.QueryTable("user_action").
		//		Filter("users_id",userInfo.Id).Filter("actions_id",actionInfo.Id).One(&userAction)
		//	if userAction.Id>0{
		//		if userAction.IsPass==1{
		//			return
		//		}else{
		//			ctx.Redirect(302,"/")
		//		}
		//	}else{//不在用户角色的情况下判断，所以下面不判断ispass ****
		//		//用户--角色--权限
		//		orm.LoadRelated(&userInfo,"roles")
		//		var actionFinal []*models.ActionInfo
		//		for _,role:=range userInfo.Roles{
		//			  orm.LoadRelated(role,"actions")
		//			  for _,action:=range role.Actions {
		//				  if action.Id == actionInfo.Id{
		//				  actionFinal = append(actionFinal, action)
		//				  }
		//			  }
		//		}
		//		if len(actionFinal)>1{
		//				return
		//		}else{
		//			ctx.Redirect(302,"/")
		//		}
		//	}
		//}else{
		//	//没有对应的方法
		//	ctx.Redirect(302,"/")
		//}
	}else{
		//用户没登录
		ctx.Redirect(302,"/")
	}
}
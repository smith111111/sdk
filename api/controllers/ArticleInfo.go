package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"api/models"
	"path"
	"strconv"
	"time"
	"os"
	"io/ioutil"
	"strings"
	"html"
	"regexp"
	"math"
	"api/db"
	"github.com/astaxie/beego/config"
	"fmt"
)

type ArticleInfoController struct {
	beego.Controller
}

func(this *ArticleInfoController) GetArticleInfoIndex(){
	this.TplName="ArticleInfo/index.html"
}

func(this *ArticleInfoController) GetFinaArticleInfoIndex(){
	this.TplName="ReportDay/index.html"
}

func(this *ArticleInfoController) GetUserOrderInfoIndex(){
	this.TplName="ReportInfo/index.html"
}

func(this *ArticleInfoController) GetArticelInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	channelName:=this.GetString("channelName")
	start:=(pageIndex-1)*pageSize;
	o:=orm.NewOrm()
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	var reportUsers []models.ReportUser
	o.QueryTable("report_user").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportUsers)
	count,_:=o.QueryTable("report_user").Filter("channel_name",temp).Count()
	this.Data["json"]=map[string]interface{}{"rows":reportUsers,"total":count}
	this.ServeJSON()
}

func(this *ArticleInfoController) ShowAddArticelInfo()  {
	o:=orm.NewOrm()
	var articelClass []models.ArticelClass
	o.QueryTable("articel_class").Filter("parent_id__gte",1).All(&articelClass)
	this.Data["articelClass"]=articelClass
	this.TplName="ArticleInfo/ShowAddArticelInfo.html"
}

func(this *ArticleInfoController) FileUp()  {
		f,h,err:=this.GetFile("fileUp")
		if err!=nil{
			this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件上传失败!!"}
		}else{
			//获取文件名称
			fileName:=h.Filename
			//获取扩展名
			fileExt:=path.Ext(fileName)
			if fileExt!=".jpg"||fileExt!=".png" {
				//获取上传文件的大小
				fileSize:=h.Size
				if fileSize<50000000 {
					//构建存储的目录
					dir:="./static/fileUp/"+strconv.Itoa(time.Now().Year())+"/"+time.Now().Month().String()+"/"+strconv.Itoa(time.Now().Day())+"/"
					_,err:=os.Stat(dir)
					if err!=nil{//表示没有文件目录
						os.MkdirAll(dir,os.ModePerm)
					}
					//文件重名
					newFileName:=strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+strconv.Itoa(time.Now().Hour())+strconv.Itoa(time.Now().Minute())+strconv.Itoa(time.Now().Minute())+strconv.Itoa(time.Now().Second())
					//构建完成路径
					fullDir:=dir+newFileName+fileExt
					err1:=this.SaveToFile("fileUp",fullDir)//保存文件
					if err1==nil {
						this.Data["json"]=map[string]interface{}{"flag":"ok","msg":fullDir}
					}else{
						this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件上传失败!!"}
					}

				}else{
					this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件类太大!!"}
				}
			}else{
				this.Data["json"]=map[string]interface{}{"flag":"no","msg":"文件类型错误!!"}
			}
		}
		this.ServeJSON()

		defer  f.Close()
	}

	func (this *ArticleInfoController) AddArtice()  {
		var articelInfo=models.ArticelInfo{}
		articelInfo.DelFlag=0
		articelInfo.AddDate=time.Now()
		articelInfo.ArticleContent=this.GetString("ArticleContent")
		articelInfo.Origin=this.GetString("Origin")
		articelInfo.Title=this.GetString("title")
		articelInfo.PhotoUrl=this.GetString("PhotoUrl")
		articelInfo.KeyWords=this.GetString("KeyWords")
		articelInfo.Intro=this.GetString("Intro")
		articelInfo.Author=this.GetString("Author")
		articelInfo.FullTitle=this.GetString("FullTitle")
		articelInfo.ModifyDate=time.Now()
		o:=orm.NewOrm()
		num,_:=o.Insert(&articelInfo)
		//获取类别编号
		classId,_:=this.GetInt("className")
		//查询类别的信息。
		var classInfo models.ArticelClass
		o.QueryTable("articel_class").Filter("id",classId).One(&classInfo)
		//创建M2M对象
		m2m:=o.QueryM2M(&articelInfo,"ArticelClasses")
		m2m.Add(classInfo)
		//生成静态页面
		CreateStaticPage(int(num))

		this.Data["json"]=map[string]interface{}{"flag":"ok"}
		this.ServeJSON()

}


func CreateStaticPage(aId int )  {
	//1:根据传递过来的文章编号，查询对应的文章信息
	var articelInfo models.ArticelInfo
	o:=orm.NewOrm()
	o.QueryTable("articel_info").Filter("id",aId).One(&articelInfo)
	dir:="./static/ArticelTemplate/ArticelTemplateInfo.html"
	file,err:=os.Open(dir)
	defer  file.Close()
	if err!=nil{
		beego.Info("文件打开失败")
	}else{
		//读取打开文件中的内容。
		content,_:=ioutil.ReadAll(file)
		//将读取的内容转换成文本字符串。
		articeContent:=string(content)
		articeContent=strings.Replace(articeContent,"$Title",articelInfo.Title,-1)
		articeContent=strings.Replace(articeContent,"$Origin",articelInfo.Origin,-1)
		articeContent=strings.Replace(articeContent,"$ArticleContent",articelInfo.ArticleContent,-1)
		articeContent=strings.Replace(articeContent,"$AddDate",articelInfo.AddDate.Format("2006-01-02"),-1)
		articeContent=strings.Replace(articeContent,"$articelId",strconv.Itoa(articelInfo.Id),-1)
		//创建文件夹。
		month:=articelInfo.AddDate.Month().String()
		var m int
		for i:=0;i<len(months) ; i++ {
			if months[i]==month {
				m=i;
				break;
			}
		}
		m=m+1
		var dirDict string//日也用md:="0"+strconv.Itoa(m)
		if m<10{
			md:="0"+strconv.Itoa(m)
			dirDict="./static/Articel/"+strconv.Itoa(articelInfo.AddDate.Year())+"/"+md+"/"+strconv.Itoa(articelInfo.AddDate.Day())+"/"
		}else{
			dirDict="./static/Articel/"+strconv.Itoa(articelInfo.AddDate.Year())+"/"+strconv.Itoa(m)+"/"+strconv.Itoa(articelInfo.AddDate.Day())+"/"
		}
		_,err:=os.Stat(dirDict)
		if err!=nil{
			os.MkdirAll(dirDict,os.ModePerm)
		}
		fullDir:=dirDict+strconv.Itoa(articelInfo.Id)+".html"
		if ioutil.WriteFile(fullDir,[]byte(articeContent),0644)==nil{
			beego.Info("写入文件成功")
		}

	}



}
//定义日期切片
var months = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}


//这里的审查词该公共的提取出来，用redis查询(第一次查数据库,存redis,后面查询redis)

//添加评论
func(this *ArticleInfoController) AddComment()  {
	var forWords[]models.SensitiveWord
	var msg=html.EscapeString(this.GetString("msg"))
	//进行禁用词判断。
	if CheckForWord(msg){
		this.Data["json"]=map[string]interface{}{"flag":"no","message":"输入的评论中含有审查词"}
	}else{
		if CheckModWord(msg) {
			var comment=models.ArticelComment{}
			comment.AddDate=time.Now()
			comment.Msg=msg
			comment.IsPass=0//表示允许展示评论
			aid,_:=this.GetInt("articelId")//文章编号
			//文章信息。&nbsp; &lt; &gt;
			o:=orm.NewOrm()
			var articelInfo models.ArticelInfo
			o.QueryTable("articel_info").Filter("id",aid).One(&articelInfo)
			comment.Articel=&articelInfo
			_,err:= o.Insert(&comment)
			if err==nil{
				this.Data["json"]=map[string]interface{}{"flag":"ok"}
			}else {
				this.Data["json"]=map[string]interface{}{"flag":"no"}
			}
		}else{
			//查询出所有的替换词。
			o:=orm.NewOrm()
			o.QueryTable("sensitive_word").Filter("is_mod",0).Filter("is_forbid",0).All(&forWords)
		}

	}

	this.ServeJSON()

}
//过滤禁用词,参数为用户评论信息
func CheckForWord(msg string )(bool)  {
	//1:查询表中的禁用词。
	o:=orm.NewOrm()
	var forWords[]models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_forbid",1).All(&forWords)
	//2:进行过滤
	// "词的名称1|词名称2|词名称3"
	var words[]string
	for i:=0;i<len(forWords);i++{
		words=append(words,forWords[i].WordPattern)
	}
	str:=strings.Join(words,"|")
	reg:=regexp.MustCompile(str)
	return reg.MatchString(msg)

}
//过滤审查词
func CheckModWord(msg string) (bool) {
	o:=orm.NewOrm()
	var forWords[]models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_mod",1).All(&forWords)
	//2:进行过滤
	// "词的名称1|词名称2|词名称3"
	var words[]string
	for i:=0;i<len(forWords);i++{
		words=append(words,forWords[i].WordPattern)
	}
	str:=strings.Join(words,"|")
	str=strings.Replace(str,"{2}",".{0,2}",-1)
	str=strings.Replace(str,"\\","\\\\",-1)
	reg:=regexp.MustCompile(str)
	return reg.MatchString(msg)
}
/*
//加载评论
func(this *ArticelInfoController) LoadCommentMsg()  {
	aid,_:=this.GetInt("articelId")//文章编号
	o:=orm.NewOrm()
	var comments[]models.ArticelComment
	//注意查询的条件指定的是外键
	o.QueryTable("articel_comment").Filter("articel_id",aid).All(&comments)
   this.Data["json"]=map[string]interface{}{"msg":comments}
   this.ServeJSON()
}*/
//加载评论
func(this *ArticleInfoController) LoadCommentMsg()  {
	aid,_:=this.GetInt("articelId")//文章编号
	//接收当前页码值。
	pageInex,_:=this.GetInt("pageIndex")
	//确定每页显示多少条记录
	pageSize:=3
	//计算总页数。
	//计算总页数，必须先计算总的记录数。
	o:=orm.NewOrm()
	recordCount,_:=o.QueryTable("articel_comment").Filter("articel_id",aid).Filter("is_pass",0).Count()
	pageCount:=int(math.Ceil(float64(recordCount)/float64(pageSize)))
	//对传递过来的页码值做校验
	if pageInex<1{
		pageInex=1
	}
	if pageInex>pageCount{
		pageInex=pageCount
	}
	//计算出start的取值。
	start:=(pageInex-1)*pageSize
	var comments[]models.ArticelComment
	o.QueryTable("articel_comment").Filter("articel_id",aid).Filter("is_pass",0).Limit(pageSize,start).All(&comments)

	strHtml:=CreatePageBar(pageInex,pageCount)//获取页码条
	this.Data["json"]=map[string]interface{}{"msg":comments,"pageBar":strHtml}
	this.ServeJSON()
}
//创建页码条（返回的是a标签）
func CreatePageBar(pageIndex,pageCount int) (strHtml string) {
	//判断总页数的取值。
	if pageCount==1{
		return  ""
	}
	start:=pageIndex-5;//计算的页码的起始位置。
	if start<1{
		start=1
	}
	end:=start+9//计算终止位置。
	if end>pageCount{
		end=pageCount
	}
	if pageIndex>1{
		strHtml="<a class='pageBarLink' href='/Admin/ArticelInfo/LoadCommentMsg?pageIndex="+strconv.Itoa(pageIndex-1)+"'>上一页</a>"
	}
	for i:=start;i<=end ;i++  {
		if pageIndex==i{//访问到前页码值，这时不需要加超链接
			strHtml=strHtml+strconv.Itoa(i)
		}else{
			strHtml=strHtml+"<a class='pageBarLink' href='/Admin/ArticelInfo/LoadCommentMsg?pageIndex="+strconv.Itoa(i)+"'>"+strconv.Itoa(i)+"</a>"
		}
	}

	if pageIndex<pageCount {
		strHtml=strHtml+"<a class='pageBarLink' href='/Admin/ArticelInfo/LoadCommentMsg?pageIndex="+strconv.Itoa(pageIndex+1)+"'>下一页</a>"
	}
	return

}

func(this *ArticleInfoController) ShowFinaUserorder()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	channelName:=this.GetString("channelName")
	start:=(pageIndex-1)*pageSize;
	o:=orm.NewOrm()
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	var reportUsers []models.ReportUser
	o.QueryTable("report_user").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportUsers)
	count,_:=o.QueryTable("report_user").Filter("channel_name",temp).Count()
	this.Data["json"]=map[string]interface{}{"rows":reportUsers,"total":count}
	//this.ServeJSON()
	this.TplName = "ReportDay/index.html"
}


func(this *ArticleInfoController) GetFinalArticelInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	channelName:=this.GetString("channelName")
	start:=(pageIndex-1)*pageSize;
	o:=orm.NewOrm()
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	var reportUsers []models.ReportUser
	o.QueryTable("report_user").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportUsers)
	count,_:=o.QueryTable("report_user").Filter("channel_name",temp).Count()
	this.Data["json"]=map[string]interface{}{"rows":reportUsers,"total":count}
	this.ServeJSON()
}




func(this *ArticleInfoController) GetUserOrderInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	channelName:=this.GetString("channelName")
	start:=(pageIndex-1)*pageSize;
	o:=orm.NewOrm()
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	var reportInfos []models.ReportInfo
	o.QueryTable("report_info").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportInfos)
	count,_:=o.QueryTable("report_info").Filter("channel_name",temp).Count()

	sum:=0.00
	for _,v:=range reportInfos{
		sum+=v.PayAccount
	}

	beego.Info("sum",sum)
	this.Data["json"]=map[string]interface{}{"rows":reportInfos,"total":count,"totalValue":sum}
	this.ServeJSON()
}



func(this *ArticleInfoController) GetChannelReportIndex()  {
	this.TplName="ReportChannel/index.html"
}

func(this *ArticleInfoController) GetUserOrderChanlInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	start:=(pageIndex-1)*pageSize;
	channelName:=this.GetString("channelName")
	o:=orm.NewOrm()
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	//var reportInfos []models.ReportInfo
	//o.QueryTable("report_info").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportInfos)
	//count,_:=o.QueryTable("report_info").Filter("channel_name",temp).Count()

	reportUsers:=[]models.ReportUser{}
	o.QueryTable("report_user").Filter("channel_name",temp).OrderBy("id").Limit(pageSize,start).All(&reportUsers)
	reportChanels := make([]models.ReportChanel,len(reportUsers))
	for k,v:= range reportUsers {
		totalPayNumber := 0
		o.Raw("select COUNT(*) payCount from report_info s where s.channel_name = ?",temp).QueryRow(&totalPayNumber)

		totalPayAccount := 0
		o.Raw("select SUM(s.pay_account) count from report_info s where s.channel_name = ?",temp).QueryRow(&totalPayAccount)


		reportChanels[k]=models.ReportChanel{Id:k+1,ChannelName:v.ChannelName,Username:v.Username,TotalPayAccount:totalPayAccount,TotalPayNumber:totalPayNumber,Uid:v.Ext1}

	}
	count,_:=o.QueryTable("report_user").Filter("channel_name",temp).OrderBy("id").Count()
	this.Data["json"]=map[string]interface{}{"rows":reportChanels,"total":count}
	this.ServeJSON()
}


func(this *ArticleInfoController) GetRetentionReportIndex()  {
	this.TplName="ReportRetention/index.html"
}

func(this *ArticleInfoController) GetRetentionInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	start:=(pageIndex-1)*pageSize;
	channelName:=this.GetString("channelName")
	//var articels[]models.ArticelInfo
	//o.QueryTable("articel_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&articels)
	//count,_:=o.QueryTable("articel_info").Filter("del_flag",0).Count()
	channelName= strings.Replace(channelName,"\"", "",-1);
	temp:=0
	if channelName == "新仙侠"{
		temp=1
	}else if channelName == "旧仙侠"{
		temp=2
	}else if channelName == "传奇"{
		temp=3
	}
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	//创建数据库连接
	databasename := iniconf.String("guozhanonedatabasename")
	guozhandb := db.GetGuoZhanonenewdb(databasename)
	ip := iniconf.String("guozhanoneip")
	if guozhandb == nil {
		beego.Info(ip, "连接异常")
		return
	}

	count:=0
	reportRetentioArrays:= make([]models.ReportRetentio,0)
	if (temp==0) {
		rows, err :=guozhandb.Query("select sd.role_name,sd.role_create_time,(select COUNT(*) from view_user_login s where s.user_login_time"+
			" >= '2020-07-10 00:00:00' and s.user_login_time <= '2020-07-11 00:00:00' )/(select COUNT(*) "+
			"from view_user_create) fristNumOneRate,(select COUNT(*) from view_user_login s where"+
			" s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time <= '2020-07-13 00:00:00' )/(select COUNT(*) "+
			"from view_user_create) fristNumTwoRate,(select COUNT(*) from view_user_login s where s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time"+
			" <= '2020-07-14 00:00:00' )/(select COUNT(*) from view_user_create) fristNumThreeRate,(select COUNT(*) from view_user_login s where s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time "+
			"<= '2020-07-15 00:00:00' )/(select COUNT(*) from view_user_create) fristNumFourRate,(select COUNT(*) from view_user_login s where "+
			" s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time <= '2020-07-16 00:00:00' )/(select COUNT(*) from view_user_create) fristNumFiveRate,"+
			"(select COUNT(*) from view_user_login s where s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time <= '2020-07-17 00:00:00' )/(select COUNT(*) from view_user_create) "+
			"fristNumSixRate,(select COUNT(*) from view_user_login s where s.user_login_time >= '2020-07-10 00:00:00' and s.user_login_time <= '2020-07-25 00:00:00' )/"+
			"(select COUNT(*) from view_user_create) fristNumSevenRate,(select COUNT(*) from view_user_login s where s.user_login_time >= '2020-07-10 00:00:00' and"+
			" s.user_login_time <= '2020-08-10 00:00:00' )/(select COUNT(*) from view_user_create) fristNumEightRate,(select COUNT(*) from view_user_login s where s.user_login_time >= "+
			"'2020-07-10 00:00:00' and s.user_login_time <= '2020-08-25 00:00:00' )/(select COUNT(*) from view_user_create) fristNumLightRate from view_user_create sd limit ?,?", start, pageSize)
  fmt.Println("r",rows)
		if err == nil {
			//遍历查询的结果集合
			for rows.Next() {
				var reportRetentio = models.ReportRetentio{}
				err = rows.Scan(&reportRetentio.RoleName, &reportRetentio.RoleCreateTime, &reportRetentio.FristNumOneRate, &reportRetentio.FristNumTwoRate,  &reportRetentio.FristNumThreeRate,  &reportRetentio.FristNumFourRate,
					&reportRetentio.FristNumFiveRate, &reportRetentio.FristNumSixRate,  &reportRetentio.FristNumSevenRate,  &reportRetentio.FristNumEightRate,  &reportRetentio.FristNumLightRate,)
				reportRetentioArrays=append(reportRetentioArrays,reportRetentio)
				if err!=nil {
					fmt.Println("err",err)
					return
				}
			}
		}

		rowTwos, errTwo :=guozhandb.Query("select COUNT(*) from view_user_create sd")
		if errTwo == nil {
			//遍历查询的结果集合
			if rowTwos.Next() {
				err = rowTwos.Scan(&count)
				if err!=nil {
					fmt.Println("err",err)
					return
				}
			}
		}
	}else if(temp==3){

		databasename := iniconf.String("chuanqionedatabasename")
		chuanqiDb := db.GetChuanqionenewdb(databasename)
		if chuanqiDb == nil {
			beego.Info(ip, "连接异常")
			return
		}
		rowsChuan, err :=chuanqiDb.Query("select  o.actorname,i.lastlogin,(select COUNT(*) from actorlogin s where s.lastlogin>=" +
			" '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-11 00:00:00' )/(select COUNT(*) from actors) fristNumOneRate,(select COUNT(*) " +
			"from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-13 00:00:00' )/(select COUNT(*) from actors) fristNumTwoRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin<= '2020-07-14 00:00:00' )/(select COUNT(*) from actors) fristNumThreeRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-15 00:00:00' )/(select COUNT(*) from actors) fristNumFourRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-16 00:00:00' )/(select COUNT(*) from actors) fristNumFiveRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-17 00:00:00' )/(select COUNT(*) from actors) fristNumSixRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-07-25 00:00:00' )/(select COUNT(*) from actors) fristNumSevenRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >= '2020-07-10 00:00:00' and s.lastlogin <= '2020-08-10 00:00:00' )/(select COUNT(*) from actors) fristNumEightRate," +
			"(select COUNT(*) from actorlogin s where s.lastlogin >='2020-07-10 00:00:00' and s.lastlogin <= '2020-08-25 00:00:00' )/(select COUNT(*) from actors) fristNumLightRate" +
			" from actors o,actorlogin i where o.accountname=i.account limit ?,?", start, pageSize)
		if err == nil {
			//遍历查询的结果集合
			for rowsChuan.Next() {
				var reportRetentio = models.ReportRetentio{}
				err = rowsChuan.Scan(&reportRetentio.RoleName, &reportRetentio.RoleCreateTime, &reportRetentio.FristNumOneRate, &reportRetentio.FristNumTwoRate,  &reportRetentio.FristNumThreeRate,  &reportRetentio.FristNumFourRate,
					&reportRetentio.FristNumFiveRate, &reportRetentio.FristNumSixRate,  &reportRetentio.FristNumSevenRate,  &reportRetentio.FristNumEightRate,  &reportRetentio.FristNumLightRate,)
				reportRetentioArrays=append(reportRetentioArrays,reportRetentio)
				if err!=nil {
					fmt.Println("err",err)
					return
				}
		     }
			}

		rowTwos, errTwo :=guozhandb.Query("select COUNT(*) from actors o,actorlogin i where o.accountname=i.account")
		if errTwo == nil {
			//遍历查询的结果集合
			if rowTwos.Next() {
				err = rowTwos.Scan(&count)
				if err!=nil {
					fmt.Println("err",err)
					return
				}
			}
		}
	}
	fmt.Println("reportRetentioArrays",reportRetentioArrays)
	fmt.Println("count",count)
	this.Data["json"]=map[string]interface{}{"rows":reportRetentioArrays,"total":count}
	this.ServeJSON()
}


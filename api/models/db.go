package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func init() {
	var dbhost string
	var dbport string
	var dbuser string
	var dbpassword string
	var db string
	//获取配置文件中对应的配置信息
	dbhost = beego.AppConfig.String("dbhost")
	dbport = beego.AppConfig.String("dbport")
	dbuser = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	db = beego.AppConfig.String("db")
	orm.RegisterDriver("mysql", orm.DRMySQL) //注册mysql Driver
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8&loc=Asia%2FShanghai"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	orm.RegisterModel(new(UserInfo), new(RoleInfo),new(ActionInfo),new(UserAction),
		new(ArticelClass),new(ArticelInfo),new(ArticelComment),new(SensitiveWord),new(ReportUser),new(ReportInfo)) //注册模型
	orm.RunSyncdb("default", false, true)
}


type UserInfo struct {
	Id        int
	UserName  string      //用户名
	UserPwd   string      //密码
	Remark    string      //备注
	AddDate   time.Time   //添加日期
	ModifDate time.Time   //修改日期
	DelFlag   int         // 删除标记 软删除。
	Roles     []*RoleInfo `orm:"rel(m2m)"`
	UserActions []*UserAction `orm:"reverse(many)"`
}

type RoleInfo struct {
	Id        int
	RoleName  string
	Remark    string
	DelFlag   int
	AddDate   time.Time
	ModifDate time.Time
	Users     []*UserInfo `orm:"reverse(many)"`
	Actions []*ActionInfo`orm:"rel(m2m)"`
}

//权限信息
type ActionInfo struct {
	Id int
	Remark string
	DelFlag int
	AddDate time.Time
	ModifDate time.Time
	Url string
	HttpMethod string
	ActionInfoName string
	ActionTypeEnum int //权限类型。
	MenuIcon string //图片地址
	IconWidth int
	IconHeight int
	Roles[]*RoleInfo `orm:"reverse(many)"`
	UserActions []*UserAction `orm:"reverse(many)"`
}
//用户权限中间表
type UserAction struct {
	Id int
	IsPass int
	Users *UserInfo `orm:"rel(fk)"`
	Actions *ActionInfo `orm:"rel(fk)"`
}

//文章类别
type ArticelClass struct {
	Id int//主键
	ClassName string//类别名称
	ParentId int //父类别的编号
	CreateUserId int //创建类别的用户编号
	CreateDate time.Time //创建时间
	DelFlag int //删除标记
	Remark string //备注
	Artices []*ArticelInfo`orm:"reverse(many)"`
}

//文章信息表
type ArticelInfo struct {
	Id int
	KeyWords string  //关键词
	Title string    //标题
	FullTitle string  //全标题
	Intro string  //导读
	ArticleContent string `orm:"type(text)"` //新闻内容  //字段较大
	Author string//作者
	Origin string//来源
	AddDate time.Time//添加日期
	ModifyDate time.Time//修改日期
	DelFlag int//删除标记
	PhotoUrl string//图片地址
	ArticelClasses []*ArticelClass`orm:"rel(m2m)"`
	Comments[]*ArticelComment `orm:"reverse(many)"`
}

//发布评论
//文章评论
type ArticelComment struct {
	Id int `from:"-"`
	Msg string
	AddDate time.Time
	IsPass int
	Articel *ArticelInfo `orm:"rel(fk)"`
}
//敏感词库
type SensitiveWord struct {
	Id int
	WordPattern string
	IsForbid int
	IsMod int
	ReplaceWord string `orm:"null"`
}

type Userorder struct {
	Id            int64  `orm:"pk;auto"`
	Channel       string `orm:"size(255)"`
	RegionalCloth string `orm:"size(255)"`
	OrderNumber   string `orm:"size(255)"`
	UserId        int64  `orm:"size(20)"`
	Role          string `orm:"size(255)"`
	AmountMoney   int
	OrderDate     time.Time `orm:"type(datetime);auto_now_add;size(20)"`
	PayStatus     int `orm:"size(11)"`
}


type PlayRole struct {
	Chanel string
	RCloth string
	OrderId string
	PlayerId int64
	ReturnTime time.Time
	ChargeStatus int
	Amount int
	Name string
}



type ReportUser struct {
	Id            int  `orm:"pk;auto"`
	Username      string `orm:"size(200)"`
	PayAccount    float64 `orm:"size(11)"`
	MonthTotal    float64 `orm:"size(11)"`
	IsActive      string  `orm:"size(200)"`
	Arppu   string  `orm:"size(200)"`
	Arpu     string  `orm:"size(200)"`
	CreateTime time.Time
	ChannelName string `orm:"size(200)"`
	Area string `orm:"size(200)"`
	Ext1 string `orm:"size(200)"`
	Ext2 string `orm:"size(200)"`
	Newuserpayrate string `orm:"size(200)"`
	Newpayuser string `orm:"size(200)"`
	Activepayrate string `orm:"size(200)"`
}



type ReportInfo struct {
	Id            int  `orm:"pk;auto"`
	Username      string `orm:"size(255)"`
	Uid          string `orm:"size(200)"`
	Area          string `orm:size(200)`
	OrderId       string `orm:"size(200)"`
	Rolename     string `orm:"size(255)"`
	PayAccount    float64  `orm:"size(11)"`
	OrderDate     string `orm:"size(255)"`
	PayStatus    int `orm:"size(4)"`
	ChannelName    string `orm:size(200)`
	Ext1 string `orm:"size(200)"`
	Ext2 string `orm:"size(200)"`
}


type ReportChanel struct {
	Id            int  `orm:"pk;auto"`//排序序号
	ChannelName    string `orm:size(200)`//默认全部渠道/可筛选渠道
	Username      string `orm:"size(255)"`//新增账号
	TotalPayNumber   int `orm:"size(11)"`//总付费额
	TotalPayAccount    int `orm:"size(11)"`//总付费人数
	Uid              string `orm:"size(200)"`//关联ID
}


type ReportRetentio struct {
	Id            int  `orm:"pk;auto"`//排序序号
	RoleCreateTime time.Time
	RoleName      string `orm:"size(255)"`//新增账号
	FristNumOneRate   string `orm:"size(200)"`
	FristNumTwoRate   string `orm:"size(200)"`
	FristNumThreeRate string `orm:"size(200)"`
	FristNumFourRate  string `orm:"size(200)"`
	FristNumFiveRate  string `orm:"size(200)"`
	FristNumSixRate   string `orm:"size(200)"`
	FristNumSevenRate string `orm:"size(200)"`
	FristNumEightRate string `orm:"size(200)"`
	FristNumLightRate string `orm:"size(200)"`
}
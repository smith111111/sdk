package utils

import (
	"api/models"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	"strconv"
	"sync"
	"time"
	"api/db"
)

var lockxianxiaone sync.Mutex
var lockguozhanone sync.Mutex
var locknewxianxiaone sync.Mutex
var lockchuanqione sync.Mutex
var secondTime = 10000 * time.Second


//同步
func InitGuozhanOneDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(secondTime) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}

	projectcode := iniconf.String("guozhanoneprojectcode")
	//创建数据库连接
	databasename := iniconf.String("guozhanonedatabasename")
	guozhandb := db.GetGuoZhanonenewdb(databasename)
	ip := iniconf.String("guozhanoneip")
	if guozhandb == nil {
		beego.Info(ip, "连接异常")
		return
	}
	for {
		select {
		case <-ticker.C:
			//	rows, err := db.Query("select r.server_id,r.role_id,r.coin,r.purchase_time,r.id from t_purchase_record r where r.pay_type = 2")
			//select user_id,SUM(role_pay_coin),role_name,user_game_server from view_user_pay GROUP BY role_name
			rows, err := guozhandb.Query("select user_game_server,role_name,SUM(role_pay_coin),role_purchase_time,role_id,user_id from view_user_pay GROUP BY role_name")
			fmt.Println("guozhan===",err)
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					lockguozhanone.Lock()
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser
					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Username, &reportUserMide.PayAccount, &reportUserMide.CreateTime, &reportUserMide.Id, &reportUserMide.Ext1)
					if err != nil {
						fmt.Println(err)
						continue
					}
					//用户不可能同区
					//reportUserMideId := strconv.Itoa(reportUserMide.Id)
					row, err := guozhandb.Query("select SUM(role_pay_coin) as month_total from view_user_pay")
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.MonthTotal)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}else {
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					//是否活跃
					count := 0
					rowCount, err := guozhandb.Query("select COUNT(*) from view_user_login where user_id = ? having count(*) >2", reportUserMide.Ext1)
					if err == nil {
						if rowCount.Next() {
							err = rowCount.Scan(&count)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}else {
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					if count == 0 {
						reportUserMide.IsActive = "0"
					} else {
						reportUserMide.IsActive = "1"
					}
					fmt.Println("13--------",reportUserMide.Username,"aa====",reportUserMide.Area)

					//orm.QueryTable("report_user").Filter("username",reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("channel_name", 0).One(&reportUser)
					fmt.Println("12--------",reportUser)
					reportUserMideIdStr := strconv.Itoa(reportUserMide.Id)

					tempStr := "4000" + reportUserMideIdStr

					tempStrInt, err := strconv.Atoi(tempStr)
					if err != nil {
						beego.Info(err)
						continue
					}
					orm.QueryTable("report_user").Filter("id",tempStrInt).One(&reportUser)
					reportUserMide.Id = tempStrInt
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount,
						ChannelName: projectcode, Id: tempStrInt,
						MonthTotal: reportUserMide.MonthTotal, CreateTime: reportUserMide.CreateTime, IsActive: reportUserMide.IsActive, Ext1: reportUserMide.Ext1}
						fmt.Println("1--------",reportUserStr)
					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr,"username","area","pay_account","channel_name","id","month_total","create_time","is_active","ext1")
						if err != nil {
							beego.Info(err)
							continue
						}
					}
					lockguozhanone.Unlock()
				}
			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
		defer	guozhandb.Close()
		}
	}

}

//旧仙侠一区
func InitXianXiaOneDataDetailsTask() {
	beego.Info("旧仙侠=======")
	orm := orm.NewOrm()
	ticker := time.NewTicker(secondTime) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	projectcode := iniconf.String("xianxiaoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	databasename := iniconf.String("xianxiaonedatabasename")
	xianxiaonedb := db.GetXianxiaonedb(databasename)
	if xianxiaonedb == nil {
		beego.Info(ip, "连接异常")
		return
	}
	timeNum := 0
	for {
		select {
		case <-ticker.C:
			rows, err := xianxiaonedb.Query("select o.create_time,o.uid,o.server_id,o.id,o.name,(o.create_time > (unix_timestamp(now()) - 2592000)) isActive ," +
				"(select IFNULL(sum(r.amount),'0') AS pay_account from player_recharge_info r where  r.player_id = o.id ) pay_account from player_info o")
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					lockxianxiaone.Lock()
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					//	err = rows.Scan(&reportUserMide.Ext1,&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					err = rows.Scan(&timeNum, &reportUserMide.Ext1, &reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.PayAccount)
					beego.Info("旧仙侠=======",err)
					if err != nil {
						fmt.Println(err)
						continue
					}
					row, err := xianxiaonedb.Query("select sum(amount) as pay_account from player_recharge_info")
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.MonthTotal)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}
					beego.Info("reportUserMide.MonthTotal",reportUserMide.MonthTotal)
					reportUserMideIdStr := strconv.Itoa(reportUserMide.Id)

					tempStr := "2" + reportUserMideIdStr

					tempStrInt, err := strconv.Atoi(tempStr)

					if err != nil {
						beego.Info(err)
						continue
					}
					tm := time.Unix(int64(timeNum), 0)
					//orm.QueryTable("report_user").Filter("username", reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("channel_name", 2).One(&reportUser)
					orm.QueryTable("report_user").Filter("id", tempStrInt).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount,
						ChannelName: projectcode, Id: tempStrInt,
						MonthTotal: reportUserMide.MonthTotal, CreateTime: tm, IsActive: reportUserMide.IsActive, Ext1: reportUserMide.Ext1}
					fmt.Println("1--------",reportUserStr,"22",reportUser.Id)
					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr,"username","area","pay_account","channel_name","id","month_total","create_time","is_active","ext1")
						if err != nil {
							beego.Info(err)
							continue
						}
					}
					lockxianxiaone.Unlock()
				}
			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer xianxiaonedb.Close()
		}
	}

}

func InitXianXiaTwoDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(3600 * time.Second) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	databasename := iniconf.String("xianxiatwodatabasename")
	username := iniconf.String("xianxiaoneusername")
	password := iniconf.String("xianxiaonepassword")
	projectcode := iniconf.String("xianxiaoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	sqlStr := username + ":" + password + "@tcp(" + ip + ")/" + databasename + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true" //todo注册多个数据库
	//创建数据库连接
	db, err := sql.Open("mysql", sqlStr)
	//关闭连接
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := db.Query("select s.server_id,s.id,s.name as username ,s.isActive as is_active ,v.totalAmount/v.`COUNT(*)` as arppu,v.totalAmount as month_total," +
				" v.totalAmount/v.`COUNT(*)` as arpu from jx_1_game_player_info_view s , jx_1_game_player_recharge_info_view v")
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					if err != nil {
						fmt.Println(err)
						return
					}
					reportUserMideId := strconv.Itoa(reportUserMide.Id)
					row, err := db.Query("select sum(amount) as pay_account from player_recharge_info where player_id = ?", reportUserMideId)
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.PayAccount)
						}
					}
					orm.QueryTable("report_user").Filter("username",reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("2",reportUserMide.ChannelName).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount, ChannelName: projectcode, Id: reportUserMide.Id, Arpu: reportUserMide.Arpu, Arppu: reportUserMide.Arppu, MonthTotal: reportUserMide.MonthTotal, CreateTime: time.Now(), IsActive: reportUserMide.IsActive}

					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					}
				}
			}
		}
	}
}

func InitXianXiaThreeDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(3600 * time.Second) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	databasename := iniconf.String("xianxiaothreedatabasename")
	username := iniconf.String("xianxiaoneusername")
	password := iniconf.String("xianxiaonepassword")
	projectcode := iniconf.String("xianxiaoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	sqlStr := username + ":" + password + "@tcp(" + ip + ")/" + databasename + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true" //todo注册多个数据库
	//创建数据库连接
	db, err := sql.Open("mysql", sqlStr)
	//关闭连接
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := db.Query("select s.server_id,s.id,s.name as username ,s.isActive as is_active ,v.totalAmount/v.`COUNT(*)` as arppu,v.totalAmount as month_total," +
				" v.totalAmount/v.`COUNT(*)` as arpu from jx_1_game_player_info_view s , jx_1_game_player_recharge_info_view v")
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					if err != nil {
						fmt.Println(err)
						return
					}
					reportUserMideId := strconv.Itoa(reportUserMide.Id)
					row, err := db.Query("select sum(amount) as pay_account from player_recharge_info where player_id = ?", reportUserMideId)
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.PayAccount)
						}
					}
					orm.QueryTable("report_user").Filter("username",reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("2",reportUserMide.ChannelName).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount, ChannelName: projectcode, Id: reportUserMide.Id, Arpu: reportUserMide.Arpu, Arppu: reportUserMide.Arppu, MonthTotal: reportUserMide.MonthTotal, CreateTime: time.Now(), IsActive: reportUserMide.IsActive}

					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					}
				}
			}
		}
	}
}

func InitXianXiaFourDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(3600 * time.Second) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	databasename := iniconf.String("xianxiafourdatabasename")
	username := iniconf.String("xianxiaoneusername")
	password := iniconf.String("xianxiaonepassword")
	projectcode := iniconf.String("xianxiaoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	sqlStr := username + ":" + password + "@tcp(" + ip + ")/" + databasename + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true" //todo注册多个数据库
	//创建数据库连接
	db, err := sql.Open("mysql", sqlStr)
	//关闭连接
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := db.Query("select s.server_id,s.id,s.name as username ,s.isActive as is_active ,v.totalAmount/v.`COUNT(*)` as arppu,v.totalAmount as month_total," +
				" v.totalAmount/v.`COUNT(*)` as arpu from jx_1_game_player_info_view s , jx_1_game_player_recharge_info_view v")
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					if err != nil {
						fmt.Println(err)
						return
					}
					reportUserMideId := strconv.Itoa(reportUserMide.Id)
					row, err := db.Query("select sum(amount) as pay_account from player_recharge_info where player_id = ?", reportUserMideId)
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.PayAccount)
						}
					}
					orm.QueryTable("report_user").Filter("username",reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("2",reportUserMide.ChannelName).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount, ChannelName: projectcode, Id: reportUserMide.Id, Arpu: reportUserMide.Arpu, Arppu: reportUserMide.Arppu, MonthTotal: reportUserMide.MonthTotal, CreateTime: time.Now(), IsActive: reportUserMide.IsActive}

					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					}
				}
			}
		}
	}
}

func InitXianXiaNewOneDataDetailsTask() {
	beego.Info("新仙侠=======")
	orm := orm.NewOrm()
	ticker := time.NewTicker(secondTime) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	projectcode := iniconf.String("xianxianewoneprojectcode")
	ip := iniconf.String("xianxianewoneip")
	databasename := iniconf.String("xianxianewonedatabasename")
	xianxiaonenewdb := db.GetXianxiaonenewdb(databasename)
	if xianxiaonenewdb == nil {
		beego.Info(ip, "连接异常")
		return
	}
	for {
		select {
		case <-ticker.C:
			beego.Info("22=======")
			rows, err := xianxiaonenewdb.Query("select o.create_time,o.uid,o.server_id,o.id,o.name,(o.create_time > (unix_timestamp(now()) - 2592000)) isActive ," +
				"(select IFNULL(sum(r.amount),'0') AS pay_account from player_recharge_info r where  r.player_id = o.id ) pay_account from player_info o;")
			timeNum := 0
			fmt.Println("新游戏",err)
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					//创建数据库连接
					locknewxianxiaone.Lock()
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					//	err = rows.Scan(&reportUserMide.Ext1,&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					err = rows.Scan(&timeNum, &reportUserMide.Ext1, &reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.PayAccount)
					if err != nil {
						fmt.Println(err)
						continue
					}
					row, err := xianxiaonenewdb.Query("select sum(amount) as pay_account from player_recharge_info")
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.MonthTotal)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}else{
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					tm := time.Unix(int64(timeNum), 0)

					reportUserMideIdStr := strconv.Itoa(reportUserMide.Id)

					tempStr := "1" + reportUserMideIdStr

					tempStrInt, err := strconv.Atoi(tempStr)
					if err != nil {
						beego.Info(err)
						continue
					}
					reportUserMide.Id = tempStrInt
					//orm.QueryTable("report_user").Filter("username", reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("channel_name", 1).One(&reportUser)
					orm.QueryTable("report_user").Filter("id", tempStrInt).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount,
						ChannelName: projectcode, Id: tempStrInt,
						MonthTotal: reportUserMide.MonthTotal, CreateTime: tm, IsActive: reportUserMide.IsActive, Ext1: reportUserMide.Ext1}
						fmt.Println("re",reportUser,"3",reportUserMide)
					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr,"username","area","pay_account","channel_name","id","month_total","create_time","is_active","ext1")
						if err != nil {
							beego.Info(err)
							continue
						}
					}
					locknewxianxiaone.Unlock()
				}

			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer xianxiaonenewdb.Close()

		}
	}
}

func InitXianXiaNewTwoDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(3600 * time.Second) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	databasename := iniconf.String("xianxianewtwodatabasename")
	username := iniconf.String("xianxianewoneusername")
	password := iniconf.String("xianxianewonepassword")
	projectcode := iniconf.String("xianxianewoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	sqlStr := username + ":" + password + "@tcp(" + ip + ")/" + databasename + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true" //todo注册多个数据库
	//创建数据库连接
	db, err := sql.Open("mysql", sqlStr)
	//关闭连接
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := db.Query("select s.server_id,s.id,s.name as username ,s.isActive as is_active ,v.totalAmount/v.`COUNT(*)` as arppu,v.totalAmount as month_total," +
				" v.totalAmount/v.`COUNT(*)` as arpu from jx_1_game_player_info_view s , jx_1_game_player_recharge_info_view v")
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser

					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Id, &reportUserMide.Username, &reportUserMide.IsActive, &reportUserMide.Arppu, &reportUserMide.MonthTotal, &reportUserMide.Arpu)
					if err != nil {
						fmt.Println(err)
						return
					}
					reportUserMideId := strconv.Itoa(reportUserMide.Id)
					row, err := db.Query("select sum(amount) as pay_account from player_recharge_info where player_id = ?", reportUserMideId)
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.PayAccount)
						}
					}
					orm.QueryTable("report_user").Filter("username",reportUserMide.Username).Filter("area",reportUserMide.Area).Filter("1",reportUserMide.ChannelName).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area, PayAccount: reportUserMide.PayAccount, ChannelName: projectcode, Id: reportUserMide.Id, Arpu: reportUserMide.Arpu, Arppu: reportUserMide.Arppu, MonthTotal: reportUserMide.MonthTotal, CreateTime: time.Now(), IsActive: reportUserMide.IsActive}

					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							return
						}
					}
				}
			}
		}
	}
}

func InitChuanqiOneDataDetailsTask() {
	orm := orm.NewOrm()
	ticker := time.NewTicker(secondTime) //每1个小时更新
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}

	projectcode := iniconf.String("chuanqioneprojectcode")
	ip := iniconf.String("chuanqioneip")
	databasename := iniconf.String("chuanqionedatabasename")
	chuanqiDb := db.GetChuanqionenewdb(databasename)
	if chuanqiDb == nil {
		beego.Info(ip, "连接异常")
		return
	}
	for {
		select {
		case <-ticker.C:
			rows, err := chuanqiDb.Query("select s.serverindex,s.actorname,s.accountname,( select IFNULL(SUM(money),'0') AS totalAmount from paylog o where o.actor_id = s.actorid ) pay_account," +
				"s.createtime,s.actorid,s.accountid from actors s;")
			beego.Info("chuanqi",err)
			floatTemp := ""
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					lockchuanqione.Lock()
					var reportUserMide models.ReportUser
					var reportUser models.ReportUser
					err = rows.Scan(&reportUserMide.Area, &reportUserMide.Username, &reportUserMide.Ext2, &floatTemp, &reportUserMide.CreateTime, &reportUserMide.Ext1, &reportUserMide.Id)
					if err != nil {
						fmt.Println(err)
						continue
					}
					floatTempStr, err := strconv.ParseFloat(floatTemp, 64)
					if err != nil {
						beego.Info(err)
						continue
					}
					reportUserMide.PayAccount = floatTempStr
                   fmt.Println("floatTempStr============",floatTempStr)
					row, err := chuanqiDb.Query("select SUM(money) as month_total from paylog")
					if err == nil {
						if row.Next() {
							err = row.Scan(&reportUserMide.MonthTotal)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}else {
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					reportUserMideIdStr := strconv.Itoa(reportUserMide.Id)

					tempStr := "300000" + reportUserMideIdStr

					tempStrInt, err := strconv.Atoi(tempStr)
					if err != nil {
						beego.Info(err)
						continue
					}
					reportUserMide.Id = tempStrInt

					//是否活跃
					count := 0
					rowCounta, err := chuanqiDb.Query("select COUNT(*)  from actorlogin where account = ? having count(*) >=2", reportUserMide.Ext2)
					if err == nil {
						if rowCounta.Next() {
							err = rowCounta.Scan(&count)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
					}else {
						if err != nil {
							beego.Info(err)
							continue
						}
					}
					if count == 0 {
						reportUserMide.IsActive = "0"
					} else {
						reportUserMide.IsActive = "1"
					}

					orm.QueryTable("report_user").Filter("id", tempStrInt).One(&reportUser)
					reportUserStr := models.ReportUser{Username: reportUserMide.Username, Area: reportUserMide.Area,
						PayAccount: reportUserMide.PayAccount, ChannelName: projectcode, Id: reportUserMide.Id,
						MonthTotal: reportUserMide.MonthTotal, CreateTime: reportUserMide.CreateTime,
						IsActive: reportUserMide.IsActive, Ext1: reportUserMide.Ext1}
                 fmt.Println("23",reportUserStr,"e",reportUser)
					if reportUser.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						fmt.Println("1212")
						_, err := orm.Update(&reportUserStr,"id","username","area","pay_account","channel_name","month_total","create_time","is_active","ext1")
						if err != nil {
							beego.Info(err)
							continue
						}
					}
					lockchuanqione.Unlock()
				}
			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}

			defer chuanqiDb.Close()
		}
	}
}



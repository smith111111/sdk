package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego"
	"time"
	"github.com/astaxie/beego/orm"
	"api/models"
	"fmt"
	"strconv"
	"sync"
	"api/db"
)

var locknewxianxiaOrderone sync.Mutex
var lockxianxiaOrderone sync.Mutex
var lockguozhanOrderone sync.Mutex
var lockchuanqiOrderone sync.Mutex

var secondTimeOrder = 10000 * time.Second


func InitXianXiaNewOneOrderDataDetailsTask(){
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	ticker := time.NewTicker(secondTimeOrder) //每1个小时更新
	projectcode := iniconf.String("xianxianewoneprojectcode")
	ip := iniconf.String("xianxianewoneip")
	orm := orm.NewOrm()
	databasename := iniconf.String("xianxianewonedatabasename")
	xianxiaoneneworderdb:=db.GetXianxiaonenewdb(databasename)
	if xianxiaoneneworderdb == nil {
		beego.Info(ip,"连接异常")
		return
	}

	for {
		select {
		case <-ticker.C:

			rows, err := xianxiaoneneworderdb.Query("select IFNULL((select sum(r.amount) from player_info o  where r.player_id = o.id ),0) pay_account," +
				"IFNULL((select o.name from player_info o  where r.player_id = o.id ),'') name,r.id,IFNULL((select o.uid from player_info o  where r.player_id = o.id ),0) player_id," +
				"r.order_id,r.notify_time,r.charge_status,IFNULL((select o.server_id from player_info o where  r.player_id = o.id ),0) area," +
				"IFNULL((select o.role_type from player_info o where  r.player_id = o.id ),0)roleType from player_recharge_info r  GROUP BY name")
			roleType:=0
			idValue:=0
			timeNum:=""
			pay:=""
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					//创建数据库连接
					locknewxianxiaOrderone.Lock()
					var reportInfoMide models.ReportInfo
					var reportInfo models.ReportInfo
					err = rows.Scan(&pay,&reportInfoMide.Username, &idValue, &reportInfoMide.Uid, &reportInfoMide.OrderId,&timeNum,&reportInfoMide.PayStatus,&reportInfoMide.Area,&roleType)
					if err != nil {
						fmt.Println(err)
						continue
					}
					beego.Info("ccccccc",pay)
					payFloat, _ := strconv.ParseFloat(pay, 64)
					reportInfoMide.PayAccount = payFloat
                if roleType == 1{
					 reportInfoMide.Rolename="呤风"
				}else if roleType == 2 {
					reportInfoMide.Rolename="月轮"
				}else if roleType == 3 {
					reportInfoMide.Rolename="期待"
				}else if roleType == 4 {
					reportInfoMide.Rolename="剑宗"
				}else{
					reportInfoMide.Rolename="无名"
				}

				 reportUserMideIdStr:=strconv.Itoa(idValue)

					tempStr:= "1" + reportUserMideIdStr

					fmt.Println("1========",tempStr)
					tempStrInt,_:=strconv.Atoi(tempStr)

					fmt.Println("2",tempStrInt)
				//	temp:=fibonacci(tempStrInt)
					reportInfoMide.Id = tempStrInt
					orm.QueryTable("report_info").Filter("id", tempStrInt).One(&reportInfo)
					reportUserStr := models.ReportInfo{
						Username:    reportInfoMide.Username, Area: reportInfoMide.Area, PayAccount: reportInfoMide.PayAccount,
						ChannelName: projectcode, Id: reportInfoMide.Id,
						OrderDate:   timeNum, Ext1: reportInfoMide.Ext1, Ext2: reportInfoMide.Ext2, PayStatus: reportInfoMide.PayStatus,
						Rolename:    reportInfoMide.Rolename, OrderId: reportInfoMide.OrderId, Uid: reportInfoMide.Uid,
					}
					if reportInfo.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					locknewxianxiaOrderone.Unlock()

				}

			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer xianxiaoneneworderdb.Close()
		}
	}

}


func InitXianXiaOneOrderDataDetailsTask(){
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	ticker := time.NewTicker(secondTimeOrder) //每1个小时更新
	projectcode := iniconf.String("xianxiaoneprojectcode")
	ip := iniconf.String("xianxiaoneip")
	orm := orm.NewOrm()
	databasename := iniconf.String("xianxiaonedatabasename")
	xianxiaoneorderdb:=db.GetXianxiaonedb(databasename)
	if xianxiaoneorderdb == nil {
		beego.Info(ip,"连接异常")
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := xianxiaoneorderdb.Query("select IFNULL((select sum(r.amount) from player_info o  where r.player_id = o.id ),0) pay_account," +
				"IFNULL((select o.name from player_info o  where r.player_id = o.id ),'') name,r.id,IFNULL((select o.uid from player_info o  where r.player_id = o.id ),0) player_id," +
				"r.order_id,r.notify_time,r.charge_status,IFNULL((select o.server_id from player_info o where  r.player_id = o.id ),0) area," +
				"IFNULL((select o.role_type from player_info o where  r.player_id = o.id ),0)roleType from player_recharge_info r  GROUP BY name")
			roleType:=0
			idValue:=0
			timeNum:=""
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					//创建数据库连接
					lockxianxiaOrderone.Lock()
					var reportInfoMide models.ReportInfo
					var reportInfo models.ReportInfo
					err = rows.Scan(&reportInfoMide.PayAccount,&reportInfoMide.Username, &idValue, &reportInfoMide.Uid, &reportInfoMide.OrderId,&timeNum,&reportInfoMide.PayStatus,&reportInfoMide.Area,&roleType)
					if err != nil {
						fmt.Println(err)
						continue
					}

					if roleType == 1{
						reportInfoMide.Rolename="呤风"
					}else if roleType == 2 {
						reportInfoMide.Rolename="月轮"
					}else if roleType == 3 {
						reportInfoMide.Rolename="期待"
					}else if roleType == 4 {
						reportInfoMide.Rolename="剑宗"
					}else{
						reportInfoMide.Rolename="无名"
					}
					reportUserMideIdStr:=strconv.Itoa(idValue)
					tempStr:= "2" + reportUserMideIdStr
					tempStrInt,_:=strconv.Atoi(tempStr)
					//	temp:=fibonacci(tempStrInt)
					reportInfoMide.Id = tempStrInt
					orm.QueryTable("report_info").Filter("id", tempStrInt).One(&reportInfo)
					reportUserStr := models.ReportInfo{
						Username:    reportInfoMide.Username, Area: reportInfoMide.Area, PayAccount: reportInfoMide.PayAccount,
						ChannelName: projectcode, Id: reportInfoMide.Id,
						OrderDate:   timeNum, Ext1: reportInfoMide.Ext1, Ext2: reportInfoMide.Ext2, PayStatus: reportInfoMide.PayStatus,
						Rolename:    reportInfoMide.Rolename, OrderId: reportInfoMide.OrderId, Uid: reportInfoMide.Uid,
					}
					if reportInfo.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					lockxianxiaOrderone.Unlock()

				}

			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer xianxiaoneorderdb.Close()
		}
	}

}


func InitGuozhanOneOrderDataDetailsTask(){
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	ticker := time.NewTicker(secondTimeOrder) //每1个小时更新
	projectcode := iniconf.String("guozhanoneprojectcode")
	ip := iniconf.String("guozhanoneip")
	orm := orm.NewOrm()
	databasename := iniconf.String("guozhanonedatabasename")
	guozhanoneorderdb:=db.GetGuoZhanonenewdb(databasename)
	if guozhanoneorderdb == nil {
		beego.Info(ip,"连接异常")
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := guozhanoneorderdb.Query("select SUM(role_pay_coin) pay_account,role_name,role_id,user_id," +
				"role_order_id,role_purchase_time,role_pay_status,user_game_server  from view_user_pay GROUP BY role_name")
			idValue:=0
			timeNum:=""
			uidstr:=""
			if err == nil {
				//遍历查询的结果集合
				for rows.Next() {
					//创建数据库连接
					lockguozhanOrderone.Lock()
					var reportInfoMide models.ReportInfo
					var reportInfo models.ReportInfo
					err = rows.Scan(&reportInfoMide.PayAccount,&reportInfoMide.Username, &idValue, &uidstr, &reportInfoMide.OrderId,&timeNum,&reportInfoMide.PayStatus,&reportInfoMide.Area)
					if err != nil {
						fmt.Println(err)
						continue
					}
					reportInfoMide.Uid=uidstr
					reportUserMideIdStr:=strconv.Itoa(idValue)
					tempStr:= "4000" + reportUserMideIdStr
					tempStrInt,_:=strconv.Atoi(tempStr)
					//	temp:=fibonacci(tempStrInt)
					reportInfoMide.Id = tempStrInt
					orm.QueryTable("report_info").Filter("id",tempStrInt).One(&reportInfo)
					reportUserStr := models.ReportInfo{
						Username:    reportInfoMide.Username, Area: reportInfoMide.Area, PayAccount: reportInfoMide.PayAccount,
						ChannelName: projectcode, Id: reportInfoMide.Id,
						OrderDate:   timeNum, Ext1: reportInfoMide.Ext1, Ext2: reportInfoMide.Ext2, PayStatus: reportInfoMide.PayStatus,
						Rolename:    reportInfoMide.Rolename, OrderId: reportInfoMide.OrderId, Uid: reportInfoMide.Uid,
					}
					if reportInfo.Id == 0 {
						_, err := orm.Insert(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					} else {
						_, err := orm.Update(&reportUserStr)
						if err != nil {
							beego.Info(err)
							continue
						}
					}

					lockguozhanOrderone.Unlock()

				}

			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer guozhanoneorderdb.Close()
		}
	}

}


func InitChuanqiOneOrderDataDetailsTask(){
	iniconf, err := config.NewConfig("ini", "./conf/mysql.conf")
	if err != nil {
		beego.Info(err)
		return
	}
	ticker := time.NewTicker(secondTimeOrder) //每1个小时更新
	projectcode := iniconf.String("chuanqioneprojectcode")
	ip := iniconf.String("chuanqioneip")
	orm := orm.NewOrm()
	databasename := iniconf.String("chuanqionedatabasename")
	chuanqioneorderdb:=db.GetChuanqionenewdb(databasename)
	if chuanqioneorderdb == nil {
		beego.Info(ip,"连接异常")
		return
	}

	for {
		select {
		case <-ticker.C:
			rows, err := chuanqioneorderdb.Query(" select DISTINCT p.rolename,p.actor_id,p.serverid from paylog p ")
			if err == nil {
					//遍历查询的结果集合
					for rows.Next() {
						//创建数据库连接
						lockchuanqiOrderone.Lock()
						var reportInfoMide models.ReportInfo
						var reportInfo models.ReportInfo
						//err = rows.Scan(&reportInfoMide.PayAccount,&reportInfoMide.Username, &idValue, &uidstr, &reportInfoMide.OrderId,&timeNum,&reportInfoMide.PayStatus,&reportInfoMide.Area)
						err = rows.Scan(&reportInfoMide.Username, &reportInfoMide.Uid, &reportInfoMide.Area)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println(" reportInfoMide.Username====================", reportInfoMide.Username)
						rowMoney, err := chuanqioneorderdb.Query(" select IFNULL(SUM(p.money),0.00) from paylog p where p.rolename = ? ", reportInfoMide.Username)
						fmt.Println("000000",err)
						if err != nil {
							fmt.Println(err)
							continue
						}
						tempMoney:=""
						if rowMoney.Next() {
							err = rowMoney.Scan(&tempMoney)

							if err != nil {
								fmt.Println(err)
								continue
							}
							fmt.Println("221212",tempMoney)

							floatTempMoney, err := strconv.ParseFloat(tempMoney, 64)
							if err != nil {
								fmt.Println(err)
								continue
							}
							reportInfoMide.PayAccount =floatTempMoney
						}

						rowOrderNo, err := chuanqioneorderdb.Query(" select max(p.order_no) from paylog p where p.actor_id = ? ", reportInfoMide.Uid)
						if err != nil {
							fmt.Println(err)
							continue
						}
						if rowOrderNo.Next() {
							err = rowOrderNo.Scan(&reportInfoMide.OrderId)
							if err != nil {

								fmt.Println(err)
								continue
							}
						}
						tempStr := "3" + reportInfoMide.Uid
						tempStrInt, err := strconv.Atoi(tempStr)

						if err != nil {
							fmt.Println(err)
							continue
						}
						//	temp:=fibonacci(tempStrInt)
						reportInfoMide.Id = tempStrInt
						orm.QueryTable("report_info").Filter("id",tempStrInt).One(&reportInfo)
						reportUserStr := models.ReportInfo{
							Username:    reportInfoMide.Username, Area: reportInfoMide.Area, PayAccount: reportInfoMide.PayAccount,
							ChannelName: projectcode, Id: reportInfoMide.Id, PayStatus: 1,
							OrderId:     reportInfoMide.OrderId, Uid: reportInfoMide.Uid,
						}
						fmt.Println("reportUserStr===",reportUserStr.PayAccount)
						if reportInfo.Id == 0 {
							_, err := orm.Insert(&reportUserStr)
							if err != nil {
								beego.Info(err)
								continue
							}
						} else {
							_, err := orm.Update(&reportUserStr)
							if err != nil {
								beego.Info(err)
								continue
							}
						}
						lockchuanqiOrderone.Unlock()
					}


			}else {
				if err != nil {
					beego.Info(err)
					continue
				}
			}
			defer chuanqioneorderdb.Close()
		}
	}

}
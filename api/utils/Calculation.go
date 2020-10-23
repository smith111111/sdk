package utils

import (
	"time"
	"github.com/astaxie/beego/orm"
	"api/models"
	"fmt"
	"github.com/astaxie/beego"
)

var secondTimeCalculation = 10000 * time.Second

func ReportUserDataCalculation() { //-select sum(pay_account) from report_user GROUP BY username

	ticker := time.NewTicker(secondTimeCalculation) //每1个小时更新

	orm := orm.NewOrm()
	reportUserArray := []models.ReportUser{}
	for {
		select {
		case <-ticker.C:
			fmt.Println("计算")
			for i := 0; i <= 3; i++ {
				orm.Raw("select id,username,month_total from report_user where channel_name = ?", i).QueryRows(&reportUserArray)
				for _, v := range reportUserArray {
					fmt.Println("reportUserArray",v.MonthTotal)
					reportUserMide := models.ReportUser{}
					aCount, _ := orm.QueryTable("report_user").Filter("channel_name", i).Filter("is_active", 1).Count()
					fmt.Println("aCount",aCount)
					bCount, _ := orm.QueryTable("report_user").Filter("channel_name", i).Count()
					fmt.Println("bCount",bCount)
					cCount := 0
					orm.Raw("select COUNT(*) from report_user  where channel_name = ? and pay_account >0", i).QueryRow(&cCount)
					fmt.Println("cCount",cCount,v.Username)
					isPayUser := ""
					dCount := 0
					orm.Raw("select COUNT(*) from report_user  where channel_name = ? and pay_account >0 and username = ?", i, v.Username).QueryRow(&dCount)
					fmt.Println("dCount======",dCount)
					if dCount == 0 {
						isPayUser = "0"
					} else {
						isPayUser = "1"
					}

					reportUserMide.Newpayuser = isPayUser

					newuserpayrate := 0.00
					if bCount != 0 {
						newuserpayrate = float64(cCount) / float64(bCount)
					}
					fmt.Println("bCount",bCount)
					reportUserMide.Newuserpayrate = fmt.Sprintf("%.4f", newuserpayrate)
					activeuserpayrate := 0.00
					if aCount != 0 {
						activeuserpayrate = float64(cCount) / float64(aCount)
					}
					reportUserMide.Activepayrate = fmt.Sprintf("%.4f", activeuserpayrate)
					if cCount != 0 {
						reportUserMide.Arppu = fmt.Sprintf("%.4f", v.MonthTotal/float64(cCount))
					} else {
						reportUserMide.Arppu = "0"
					}
					if aCount != 0 {
						reportUserMide.Arpu = fmt.Sprintf("%.4f", float64(v.MonthTotal)/float64(aCount))
					} else {
						reportUserMide.Arpu = "0"
					}
					reportUserStr := models.ReportUser{Id: v.Id,Arpu: reportUserMide.Arpu, Arppu: reportUserMide.Arppu, Newuserpayrate: reportUserMide.Newuserpayrate, Activepayrate: reportUserMide.Activepayrate, Newpayuser: isPayUser}
                     fmt.Println("2===",reportUserStr)
					_, err := orm.Update(&reportUserStr,"arpu","arppu","newuserpayrate","activepayrate","newpayuser")
					if err != nil {
						beego.Info(err)
						return
					}
				}

			}
		}
	}
}

package mysql

import (
	"akali/app/models"
	"log"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

// Seed 建立測試資料
func (db *DB) Seeds(tools *tools.Tools) {
	users := []models.User{
		{
			Name: "阿卡莉",
			Ips:  `["192.168.1.10", "10.0.0.5"]`,
		},
	}

	for _, user := range users {
		user.CreateTime = tools.NowUnix()
		user.UpdateTime = tools.NowUnix()
		var exists int64
		db.WDB.Model(&models.User{}).Where("name = ?", user.Name).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&user)
		}
	}
	log.Println("Database seed complete")
}

package attendance

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Attendance struct {
	gorm.Model
	User_id    string
	Start_time string
	End_time   string
}

func connectSql() (DB *gorm.DB, err error) {
	DBMS := "mysql"
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER_NAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PROTOCOL := "tcp(" + DBMS + ":" + os.Getenv("DB_PORT") + ")"

	dsn := DB_USER_NAME + ":" + DB_PASSWORD + "@" + DB_PROTOCOL + "/" + DB_NAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func InsertRecord(user_id string, current_time string) {
	db, err := connectSql()
	if err != nil {
		log.Fatal(err)
	}

	db.Create(&Attendance{User_id: user_id, Start_time: current_time})
}

func UpdateRecord(user_id string, current_time string) {
	db, err := connectSql()
	if err != nil {
		log.Fatal(err)
	}

	db.Model(&Attendance{}).Where("user_id = ?", user_id).Where("end_time = ?", "").Update("end_time", current_time)
}

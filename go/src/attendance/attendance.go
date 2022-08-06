package attendance

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

type Attendance struct {
	gorm.Model
	Id            int
	User_id       string
	Start_time    string
	End_time      string
	Working_time  string
	Breaking_time string
}

func ConnectSql() (DB *gorm.DB, err error) {
	DBMS := "mysql"
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER_NAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PROTOCOL := "tcp(" + DBMS + ":" + os.Getenv("DB_PORT") + ")"

	dsn := DB_USER_NAME + ":" + DB_PASSWORD + "@" + DB_PROTOCOL + "/" + DB_NAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func InsertRecord(user_id string, current_time string) {
	db, err := ConnectSql()
	if err != nil {
		log.Fatal(err)
	}

	db.Create(&Attendance{User_id: user_id, Start_time: current_time, Breaking_time: "0.00"})
}

func UpdateRecord(user_id string, current_time string) (string, string, string, string) {
	db, err := ConnectSql()
	if err != nil {
		log.Fatal(err)
	}

	var attendance Attendance
	db.Where("user_id = ?", user_id).Where("end_time = ?", "").Last(&attendance)

	attendance.End_time = current_time
	attendance.Working_time = CalcDurationTime(attendance.Start_time, attendance.End_time)
	db.Save(&attendance)

	return attendance.Start_time, attendance.End_time, attendance.Working_time, attendance.Breaking_time
}

func CalcDurationTime(start_time string, end_time string) string {
	layout := "2006-01-02 15:04"
	start_date_time, _ := time.Parse(layout, start_time)
	end_date_time, _ := time.Parse(layout, end_time)
	working_time := end_date_time.Sub(start_date_time)
	return strconv.FormatFloat(working_time.Hours(), 'f', 2, 64)
}

func GetAttendanceIdByUserId(user_id string, db *gorm.DB) int {
	var attendance Attendance
	result := db.Where("user_id = ?", user_id).Where("end_time = ?", "").Last(&attendance)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0
	}
	return attendance.Id
}

func AddBreakingTime(attendance_id int, breaking_time string) {
	db, err := ConnectSql()
	if err != nil {
		log.Fatal(err)
	}

	var attendance Attendance
	db.Where("id = ?", attendance_id).Last(&attendance)
	attendance.Breaking_time = addDuration(breaking_time, attendance.Breaking_time)
	db.Save(&attendance)
}

func addDuration(target_duration string, old_duration string) string {
	float_target_duration, _ := strconv.ParseFloat(target_duration, 64)
	float_old_duration, _ := strconv.ParseFloat(old_duration, 64)
	sum_duration := float_target_duration + float_old_duration
	return strconv.FormatFloat(sum_duration, 'f', 2, 64)
}

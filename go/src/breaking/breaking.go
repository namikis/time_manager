package breaking

import (
	// "errors"
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"mymodule/attendance"
	// "os"
	// "strconv"
	// "time"
)

type Break struct {
	gorm.Model
	Attendance_id    int
	Start_break_time string
	End_break_time   string
	Breaking_time    string
}

func InsertBreak(start_time string, user_id string) int {
	db, err := attendance.ConnectSql()
	if err != nil {
		log.Fatal(err)
	}
	attendance_id := attendance.GetAttendanceIdByUserId(user_id, db)
	if attendance_id == 0 {
		// アクティブな勤怠が存在しない
		return 0
	}

	db.Create(&Break{Attendance_id: attendance_id, Start_break_time: start_time})
	return 1
}

func UpdateBreak(end_time string, user_id string) int {
	db, err := attendance.ConnectSql()
	if err != nil {
		log.Fatal(err)
	}
	attendance_id := attendance.GetAttendanceIdByUserId(user_id, db)
	if attendance_id == 0 {
		// アクティブな勤怠が存在しない
		return 0
	}

	var target_breaking Break
	db.Where("attendance_id = ?", attendance_id).Where("end_break_time = ?", "").Last(&target_breaking)
	target_breaking.End_break_time = end_time
	target_breaking.Breaking_time = attendance.CalcDurationTime(target_breaking.Start_break_time, target_breaking.End_break_time)
	db.Save(&target_breaking)

	attendance.AddBreakingTime(attendance_id, target_breaking.Breaking_time)

	return 1
}

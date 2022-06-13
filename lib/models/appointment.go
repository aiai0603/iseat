package models

import (
	"fmt"
	"lib/utils"
	"strconv"
	"time"
)

type Appointment struct {
	Base
	StudentID 	uint		`json:"student" gorm:"not null"`
	SeatID		uint		`json:"seat" gorm:"not null"`
	Date		string		`json:"date" gorm:"not null"`
	Time		uint		`json:"time" gorm:"not null"`
	State       uint		`json:"state" gorm:"not null"` //0-预约中，1-签到，2-签到超时，3-取消
}

func (app *Appointment)CreateAppointment() error {
	return  GetDB().Create(&app).Error
}

func GetNowAppointmentList(times int) []Appointment{
	now := time.Now()
	date:= strconv.Itoa(now.Year())+"-"+strconv.Itoa(int(now.Month()))+"-"+strconv.Itoa(now.Day())
	var appointmentList []Appointment
	GetDB().Where("date = ? and time = ? and state = 0",date, times).Find(&appointmentList)
	return appointmentList
}

func GetNowUndoneAppointmentList(times int) []Appointment{
	now := time.Now()
	date:= strconv.Itoa(now.Year())+"-"+strconv.Itoa(int(now.Month()))+"-"+strconv.Itoa(now.Day())
	var appointmentList []Appointment
	GetDB().Model(&Appointment{}).Where("date = ? and time = ? and state = 0",date, times).Update("state","2")
	//GetDB().Where("date = ? and time = ? and state = 0",date, times).Find(&appointmentList)
	return appointmentList
}

func (app *StuAppointment)StudentCreateAppointment(stu Student,message *utils.Message) error {
	hadApp := Appointment{}
	re :=GetDB().Where("student_id = ? AND date = ? AND time = ? AND state = 0",stu.ID,app.Date,app.Time).First(&hadApp)
	if re.RowsAffected != 0{
		msg := message.GetMsg(utils.MsgAppointment, "had_appointment")
		return fmt.Errorf(msg)
	}
	var appSeats []uint
	err := GetDB().Model(Appointment{}).Select("seat_id").Where("date = ? and time = ? and state = 1",app.Date,app.Time).Find(&appSeats).Error
	if err != nil{
		msg := message.GetMsg(utils.MsgAppointment, "get_seats")
		return fmt.Errorf(msg)
	}
	db := GetDB()
	if len(appSeats) != 0{
		db = db.Not("id in ?",appSeats)
	}
	seat := Seat{}
	db = db.Where("area_id = ? ",app.Area).First(&seat)
	err = db.Error
	if err != nil{
		msg := message.GetMsg(utils.MsgAppointment, "no_seat")
		return fmt.Errorf(msg)
	}
	app.SeatID = seat.ID
	app.State = 0
	err = GetDB().Create(&app.Appointment).Error
	if err != nil{
		return err
	}
	return nil
}

func StudentGetAppointmentList(stu Student)([]AppointmentHis, error){
	var apps []Appointment
	var hiss []AppointmentHis
	err := GetDB().Where("student_id = ?",stu.ID).Order("created_at desc").Find(&apps).Error
	if err != nil{
		return hiss,err
	}
	for _,app := range apps{
		his := AppointmentHis{Appointment:app}
		seat := Seat{}
		GetDB().Where("id",his.Appointment.SeatID).First(&seat)
		his.Seat = seat.SeatID
		hiss = append(hiss,his)
	}
	return  hiss,nil
}

func StudentGetAppointmentInList(stu Student)([]AppointmentHis, error){
	var apps []Appointment
	var hiss []AppointmentHis
	err := GetDB().Where("student_id = ? And state = 0",stu.ID).Order("created_at desc").Find(&apps).Error
	if err != nil{
		return hiss,err
	}
	for _,app := range apps{
		his := AppointmentHis{Appointment:app}
		seat := Seat{}
		GetDB().Where("id",his.Appointment.SeatID).First(&seat)
		his.Seat = seat.SeatID
		hiss = append(hiss,his)
	}
	return  hiss,nil
}

func DeleteAppointment(appID uint) error  {
	return GetDB().Model(&Appointment{}).Where("id = ?", appID).Update("state", 3).Error
}

func GetAppointmentList(filer map[string]interface{})([]AppointmentHis,int64,error){
	var app []Appointment
	var apph []AppointmentHis
	var count  int64
	db := GetDB()
	if filer["date"] != nil && filer["date"].(string) != ""{
		db = db.Where("date = ?",filer["date"].(string))
	}
	if filer["state"] != nil{
		db =db.Where("state = ?",int(filer["state"].(float64)))
	}
	if filer["time"] != nil{
		db =db.Where("time = ?",int(filer["time"].(float64)))
	}
	db = db.Model(Appointment{}).Count(&count).Order("created_at desc")
	if err:=db.Error; err != nil{
		return apph,count, err
	}
	if filer["page"]!=nil && filer["limit"]!=nil{
		page := int(filer["page"].(float64))
		limit := int(filer["limit"].(float64))
		skip := (page-1) * limit
		if err := db.Offset(skip).Limit(limit).Find(&app).Error; err != nil {
			return apph,count, err
		}
	}
	for _,a :=  range app{
		seat := Seat{}
		student := Student{}
		GetDB().Where("id = ?",a.SeatID).First(&seat)
		GetDB().Where("id = ?",a.StudentID).First(&student)
		ah := AppointmentHis{Appointment:a,Seat: seat.SeatID,Name: student.Name}
		apph = append(apph,ah)
	}
	return apph,count, nil
}
type StuAppointment struct {
	Area	uint  `json:"area"'`
	Appointment
}

type AppointmentHis struct {
	Appointment
	Seat    string  `json:"seat_id"'`
	Name    string  `json:"name"'`
}

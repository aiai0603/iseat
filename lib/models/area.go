package models

import "fmt"

type Area struct {
	Base
	Name 		string	`json:"name"`
	FloorID 	uint	`json:"floorID"`
	Seats		[]Seat  `json:"-"`
}

func (area *Area)CreateArea() error{
	err := GetDB().Create(area).Error
	return err
}

func GetAreaByFloor(filer map[string]interface{}) ([]Area, error){
	var areas []Area
	err := GetDB().Order("created_at desc").Where("floor_id = ?",int(filer["searchParams"].(float64))).Find(&areas).Error
	return areas, err
}

func StudentGetAreaByFloor(id uint) ([]Area, error){
	var areas []Area
	err := GetDB().Order("created_at desc").Where("floor_id = ?",id).Find(&areas).Error
	return areas, err
}

func (area *Area)UpdateArea() error {
	err := GetDB().Model(Area{}).Where("id = ?",area.ID).Updates(area).Error
	return err
}

func (stuArea *StudentArea)GetAreaInf(app Appointment)error{
	var	seats []uint
	err := GetDB().Model(&Seat{}).Select("id").Where("area_id = ?",stuArea.ID).Count(&stuArea.CountSeat).Find(&seats).Error
	fmt.Println(seats)
	if err != nil {
		return err
	}
	err = GetDB().Model(&Appointment{}).Where("date = ? AND time = ? AND seat_id in ?",app.Date, app.Time, seats).Count(&stuArea.SeatAppointment).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteArea(id uint) error{
	return GetDB().Delete(&Area{},id).Error
}

type StudentArea struct {
	Area
	CountSeat  int64  `json:"count_seat"`
	SeatAppointment int64  `json:"seat_appointment"`
}
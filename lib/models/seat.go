package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Seat struct {
	Base
	SeatID  		string         `json:"seatID" gorm:"not null;unique"`
	AreaID  		uint		   `json:"area" gorm:"not null"`
	State			uint		   `json:"state" gorm:"not null;default:1"` //0-禁用，1-启用
	LoraID          string		   `json:"loraID"`
	Appointment		[]Appointment  `json:"-"`
}

func (seat *Seat)CreateSeat() error{
	err := GetDB().Create(seat).Error
	return err
}

func GetSeatByArea(filer map[string]interface{}) ([]Seat,int64, error){
	db := GetDB()
	var count int64
	var seat []Seat
	if filer["seatID"] != nil && filer["seatID"].(string) != ""{
		db =db.Where("seat_id like ?","%"+filer["seatID"].(string)+"%")
	}
	if filer["area_id"] != nil {
		db =db.Where("area_id = ?",int(filer["area_id"].(float64)))
	}
	if filer["state"] != nil {
		db =db.Where("state = ?",int(filer["state"].(float64)))
	}
	db = db.Model(Seat{}).Count(&count).Order("created_at desc")
	if err:=db.Error; err != nil{
		return seat,count, err
	}
	if filer["page"]!=nil && filer["limit"]!=nil{
		page := int(filer["page"].(float64))
		limit := int(filer["limit"].(float64))
		skip := (page-1) * limit
		if err := db.Offset(skip).Limit(limit).Find(&seat).Error; err != nil {
			return seat,count, err
		}
	}
	return seat,count, nil
}

func (seat *Seat)UpdateSeat() error {
	err := GetDB().Model(&Seat{}).Where("id = ?",seat.ID).Updates(seat).Error
	return err
}

func GetSeat(id uint) Seat{
	var seat Seat
	GetDB().Where("id = ?", id).First(&seat)
	return seat
}

func SeatState(id uint, state int) error{
	if state==0{
		GetDB().Model(Appointment{}).Where("state = 0 AND seat_id = ?",id).Update("state",3)
	}
	err := GetDB().Model(Seat{}).Where("id = ?",id).Update("state",state).Error
	return err
}

func DeleteSeat(id uint) error{
	seat:= GetSeat(id)
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("seat_id=?", id).Delete(&Appointment{}).Error
		if err != nil {
			return err
		}
		err = tx.Select(clause.Associations).Delete(&seat).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
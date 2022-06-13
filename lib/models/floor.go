package models

import (
	"fmt"
	"lib/utils"
)

type Floor struct {
	Base
	Name  string  `json:"name" gorm:"not null;unique"`
	Areas []Area  `json:"-"`
}

func (floor *Floor)CreateFloor() error{
	err := GetDB().Create(floor).Error
	return err
}

func GetAllFloors(filer map[string]interface{})([]Floor, error){
	var floors []Floor
	var err error
	if filer["name"] == "" || filer["name"] == nil{
		err = GetDB().Find(&floors).Order("id desc").Error
	}else {
		err = GetDB().Where("name like ?","%"+filer["name"].(string)+"%").Find(&floors).Order("id desc").Error
	}

	return floors, err
}

func GetFloors(id uint, message *utils.Message)(*Floor, error){
	floor := &Floor{}
	GetDB().Where("id = ?", id).First(floor)
	if floor.Name == "" {
		msg := message.GetMsg(utils.MsgFloor, "find_floor")
		return nil, fmt.Errorf(msg)
	}
	return floor, nil
}

func DeleteFloor(floor *Floor) error{
	err := GetDB().Delete(floor).Error
	return err
}

func (floor *Floor)UpdateFloor() error{
	err := GetDB().Save(floor).Error
	return err
}

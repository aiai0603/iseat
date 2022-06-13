package models

import (
	"strconv"
	"time"
)

type Pagination struct {
	Meta Metadata    `json:"pagination"`
	Data interface{} `json:"data"`
}

type Metadata struct {
	Total    int `json:"total"`
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type Table struct {
	Code int`json:"code"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
}

func GetCount()(InitRes, error){
	init := InitRes{}
	now := time.Now()
	date:= strconv.Itoa(now.Year())+"-"+strconv.Itoa(int(now.Month()))+"-"+strconv.Itoa(now.Day())
	init.Data = date
	if err := GetDB().Model(&Appointment{}).Where("date = ? and time = 1",date).Count(&init.Morning).Error; err != nil{
		return InitRes{}, err
	}
	if err := GetDB().Model(&Appointment{}).Where("date = ? and time = 2",date).Count(&init.Noon).Error; err != nil{
		return InitRes{}, err
	}
	if err := GetDB().Model(&Appointment{}).Where("date = ? and time = 3",date).Count(&init.Night).Error; err != nil{
		return InitRes{}, err
	}
	if err := GetDB().Model(&Seat{}).Count(&init.Seats).Error; err != nil{
		return InitRes{}, err
	}
	if err := GetDB().Model(&Student{}).Count(&init.Students).Error; err != nil{
		return InitRes{}, err
	}
	if err := GetDB().Model(&Appointment{}).Count(&init.Total).Error; err != nil{
		return InitRes{}, err
	}
	return init,nil
}

func GetWeek()(InitWeek, error){
	now := time.Now()
	initWeek := InitWeek{}
	var err error
	for i := 6;i>=0;i--{
		date:= strconv.Itoa(now.Year())+"-"+strconv.Itoa(int(now.Month()))+"-"+strconv.Itoa(now.Day()-i)
		initWeek.Week[6-i]=date
		err = GetDB().Model(&Appointment{}).Where("date = ?",date).Count(&initWeek.Today[6-i]).Error
		if err != nil{
			break
		}
		err = GetDB().Model(&Appointment{}).Where("date = ? and time=1",date).Count(&initWeek.Morning[6-i]).Error
		if err != nil{
			break
		}
		err = GetDB().Model(&Appointment{}).Where("date = ? and time=2",date).Count(&initWeek.Noon[6-i]).Error
		if err != nil{
			break
		}
		err = GetDB().Model(&Appointment{}).Where("date = ? and time=3",date).Count(&initWeek.Night[6-i]).Error
		if err != nil{
			break
		}
	}
	return initWeek,err
}

type InitRes struct {
	Morning int64	`json:"morning"`
	Noon    int64	`json:"noon"`
	Night   int64	`json:"night"`
	Seats   int64	`json:"seats"`
	Students int64	`json:"students"`
	Total    int64	`json:"total"`
	Data   string
}

type InitWeek struct {
	Week  [7]string `json:"week"`
	Today [7]int64 `json:"tokay"`
	Morning [7]int64 `json:"morning"`
	Noon [7]int64 `json:"noon"`
	Night[7]int64 `json:"night"`
}
package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lib/utils"
)

type Student struct {
	Base
	Username  string    	`json:"username" gorm:"unique;not null"`
	Password  string		`json:"password" gorm:"not null"`
	Name	  string		`json:"name" gorm:"not null"`
	State  	  int			`json:"state" gorm:"not null;;default:1"` // 0-封禁，1-正常
	Appointment []Appointment  `json:"-"`
 }

type ResStudentLogin struct {
	Student
	Token string `json:"token"`
}


func (stu *Student)CreateStudent()error{
	stu.Password = utils.Md5V2("123456")
	err := GetDB().Create(stu).Error
	return err
}


func DeleteStudent(id uint) error{
	student,_:= GetStudent(id)
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("student_id=?", id).Delete(&Appointment{}).Error
		if err != nil {
			return err
		}
		err = tx.Select(clause.Associations).Delete(&student).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetStudent(id uint) (Student,error){
	var stu Student
	err := GetDB().Where("id = ?", id).First(&stu).Error
	return stu,err
}

func GetStudentList(filer map[string]interface{}) ([]Student,int64,error) {
	var student []Student
	var count  int64
	db := GetDB()
	if filer["username"] != nil && filer["username"].(string) != ""{
		db = db.Where("username Like ?","%"+filer["username"].(string)+"%")
	}
	if filer["name"] != nil && filer["name"].(string) != ""{
		db = db.Where("name Like ?","%"+filer["name"].(string)+"%")
	}
	if filer["state"] != nil{
		db =db.Where("state = ?",int(filer["state"].(float64)))
	}
	db = db.Model(Student{}).Count(&count).Order("created_at desc")
	if err:=db.Error; err != nil{
		return student,count, err
	}
	if filer["page"]!=nil && filer["limit"]!=nil{
		page := int(filer["page"].(float64))
		limit := int(filer["limit"].(float64))
		skip := (page-1) * limit
		if err := db.Offset(skip).Limit(limit).Find(&student).Error; err != nil {
			return student,count, err
		}
	}
	return student,count, nil
}

func (student *Student)UpdateStudent() error {
	err := GetDB().Model(&Student{}).Where("id = ?",student.ID).Updates(student).Error
	return err
}

func (stu *Student)GetAppointmentCount()(int64, error){
	var count int64
	err := GetDB().Model(Appointment{}).Where("student_id = ?",stu.ID).Count(&count).Error
	return count,err
}

func (stu *Student)ChangePwd(newPwd string) error{
	pwd := utils.Md5V2(newPwd)
	err := GetDB().Model(&Student{}).Where("id = ?",stu.ID).Update("password",pwd).Error
	return err
}

func StudentState(id uint, state int) error{
	if state==0{
		GetDB().Model(Appointment{}).Where("state = 0 AND student_id = ?",id).Update("state",3)
	}
	err := GetDB().Model(Student{}).Where("id = ?",id).Update("state",state).Error
	return err
}

func UpdateStudentPWD(id uint) error{
	err := GetDB().Model(Student{}).Where("id = ?",id).Update("password",utils.Md5V2("123456")).Error
	return err
}

type StudentInf struct {
	Student
	Count  int64 `json:"count"`
}

type StudentPwd struct {
	OldPwd string `json:"oldPwd"`
	NewPwd string `json:"NewPwd"`
}
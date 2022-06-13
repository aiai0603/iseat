package models

type Notice struct {
	Base
	Title  	  string	`json:"title" gorm:"type:VARCHAR(64);unique;not null"`
	Content   string	`json:"content" gorm:"type:VARCHAR(1024);not null"`
	AdminID   uint 		`json:"adminID" gorm:"not null"`
}

func (notice *Notice)CreateNotice()error{
	err := GetDB().Create(notice).Error
	return err
}

func (notice *Notice)UpdateNotice() error{
	err := GetDB().Model(&Notice{}).Where("id = ?",notice.ID).Updates(notice).Error
	return err
}

func DeleteNotice(notice Notice) error{
	err := GetDB().Delete(&notice).Error
	return err
}

func GetNotice(id uint) (Notice, error){
	var notice Notice
	err := GetDB().Where("id = ?", id).First(&notice).Error
	return notice, err
}

func ListNotices(filer map[string]interface{}) ([]Notice,int64,error) {
	db := GetDB()
	var notice []Notice
	var count int64
	if filer["title"]!=nil && filer["title"].(string)!=""{
		db = db.Where("title LIKE ?", "%"+filer["title"].(string)+"%")
	}

	db = db.Model(Notice{}).Count(&count).Order("created_at desc")
	if err:=db.Error; err != nil{
		return notice,count, err
	}
	if filer["page"]!=nil && filer["limit"]!=nil{
		page := int(filer["page"].(float64))
		limit := int(filer["limit"].(float64))
		skip := (page-1) * limit
		if err := db.Offset(skip).Limit(limit).Find(&notice).Error; err != nil {
			return notice,count, err
		}
	}
	return notice,count, nil
}

func StudentGetNoticeList() ([]Notice,error){
	var notices []Notice
	err := GetDB().Order("created_at desc").Limit(10).Find(&notices).Error
	return notices,err
}
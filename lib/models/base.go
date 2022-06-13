package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lib/utils"

	"time"
)

var (
	db *gorm.DB
)

type Base struct {
	ID        uint  	 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	message   *utils.Message
}

func (b *Base) SetMessage(message *utils.Message) {
	b.message = message
}

func (b *Base) GetMsg(section, key string) string {
	if nil == b.message {
		message := new(utils.Message)
		return message.GetMsg(section, key)
	}
	return b.message.GetMsg(section, key)
}

//func (b *Base) BeforeCreate(db *gorm.DB) (err error) {
//	if b.ID == uuid.Nil {
//		u, err := uuid.NewRandom()
//		if err != nil {
//			return err
//		}
//		b.ID = u
//	} else { //adaptor for upsert. Only exist id can be used
//		t := db.Statement.Table
//		err = db.Table(t).First(&Base{}, "id = ?", b.ID).Error
//		if err != nil {
//			return fmt.Errorf("uuid must be generated")
//		}
//	}
//	return nil
//}

func init() {
	dsn :=	"root:Wuzufeng123!@tcp(123.207.25.242:3306)/lib?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("cannot open db, please check with the admin:" + err.Error())
	}

	if utils.GetConf().Logger.Debug {
		db = conn.Debug()
	} else {
		db = conn
	}

	//err = db.SetupJoinTable(&Device{}, "Group", &DeviceGroupRef{})
	//if err != nil {
	//	panic("cannot migrate db, please check with the admin:" + err.Error())
	//}
	err = db.AutoMigrate(&Admin{},&Floor{},&Area{},&Seat{},&Student{},&Appointment{},&Notice{})
	if err != nil {
		panic("cannot migrate db, please check with the admin:" + err.Error())
	}
	//db.Exec("ALTER table devices ADD CONSTRAINT devices_device_type_id_device_id_key UNIQUE(device_type_id, device_id);")
	//db.Exec("ALTER table device_group_refs ADD CONSTRAINT device_group_refs_device_id_device_group_id UNIQUE(device_id, device_group_id);")
	//if err := db.Exec("ALTER TABLE tags ADD UNIQUE(tag_type, name)").Error; nil != err {
	//	panic("cannot add unique to the table tag:" + err.Error())
	//}
}

func GetDB() *gorm.DB {
	return db
}

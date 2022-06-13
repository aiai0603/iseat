package utils

import (
	"io/ioutil"
	"net/http"
	"path"

	ini "gopkg.in/ini.v1"
)

const (
	MsgAdmin	  = "admin"
	MsgStudent    = "student"
	MsgFloor      = "floor"
	MsgArea       = "area"
	MsgSeat       = "seat"
	MsgToken      = "token"
	MsgNotice     = "notice"
	MsgModel      = "model"
	MsgAppointment=	"appointment"
	MsgStatistic  = "statistic"
)

var (
	g_msg map[string]*ini.File
)

type Messenger interface {
	GetMsg(section, key string) string
}

type Message struct {
	language string
}

func (msg *Message) SetLanguage(r *http.Request) {
	if nil != r {
		msg.language = r.Header.Get("Content-Language")
	}
}

func (msg *Message) GetMsg(section, key string) string {
	language := ""
	if 0 == len(msg.language) {
		language = "zh_CN.ini"
	} else {
		language = msg.language + ".ini"
	}
	if conf, ok := g_msg[language]; ok {
		s := conf.Section(section)
		if s != nil {
			return s.Key(key).String()
		}
	}
	return ""
}

func readMsgDir() error {
	g_msg = make(map[string]*ini.File)
	confDir, err := GetConfLoc()
	if nil != err {
		return err
	}

	dir := path.Join(confDir, "multilingual")
	infos, err := ioutil.ReadDir(dir)
	if nil != err {
		return err
	}

	for _, info := range infos {
		fName := info.Name()
		Log.Infof("uiMsg file : %s", fName)
		fPath := path.Join(dir, fName)
		if conf, err := ini.Load(fPath); nil != err {
			return err
		} else {
			g_msg[fName] = conf
		}
	}
	return nil
}

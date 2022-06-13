package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"lib/controllers"
	"lib/middlewares"
	"lib/mqttTask"
	"lib/utils"
	"net/http"
)
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
//func cors(f http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
//		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
//		w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
//		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
//		w.Header().Set("content-type", "application/json;charset=UTF-8")             //返回数据格式是json
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusNoContent)
//			return
//		}
//		f(w, r)
//	}
//}

func main() {
	// MQTT 初始化
	mqttTask.MqttInitialization()

	// 定时任务
	c := cron.New()
	// 预约信息
	c.AddFunc("00 8 * * *", func() {
		mqttTask.SendAppointment(0)
	})
	c.AddFunc("00 12 * * *", func() {
		mqttTask.SendAppointment(1)
	})
	c.AddFunc("00 18 * * *", func() {
		mqttTask.SendAppointment(2)
	})

	// 未签到信息
	c.AddFunc("00 9 * * *", func() {
		mqttTask.SendUndoneAppointment(0)
	})
	c.AddFunc("00 13 * * *", func() {
		mqttTask.SendUndoneAppointment(1)
	})
	c.AddFunc("00 19 * * *", func() {
		mqttTask.SendUndoneAppointment(2)
	})

	//c.AddFunc("@every 5s", func() {
	//	mqttTask.SendUndoneAppointment(1)
	//})
	c.Start()

	router := mux.NewRouter()

	// 小程序
	router.HandleFunc("/student/login", controllers.StudentLogin).Methods(http.MethodPost)  //用户登陆
	router.HandleFunc("/student/inf", controllers.StudentInf).Methods(http.MethodGet) //用户获得个人信息
	router.HandleFunc("/student/changePwd", controllers.StudentChangePwd).Methods(http.MethodPost) //用户修改密码
	router.HandleFunc("/student/token", controllers.VerifyToken).Methods(http.MethodGet) //token验证
	router.HandleFunc("/student/notices", controllers.StudentGetNotices).Methods(http.MethodGet) //获得最近10条通知
	router.HandleFunc("/student/notice/{id}", controllers.NoticeHandler).Methods(http.MethodGet) //获得通知信息
	router.HandleFunc("/student/floors", controllers.FloorHandler).Methods(http.MethodGet) //获得楼层
	router.HandleFunc("/student/area/{floor_id}", controllers.StudentGetArea).Methods(http.MethodPost) //获取楼层下区域以及座位信息
	router.HandleFunc("/student/appointment", controllers.StudentAppointmentsHandle).Methods(http.MethodPost,http.MethodGet) //预约座位，查看预约历史
	router.HandleFunc("/student/appointmentIn", controllers.StudentAppointmentIn).Methods(http.MethodGet) //获得预约中记录
	router.HandleFunc("/student/appointment/{id}", controllers.StudentDelAppointment).Methods(http.MethodDelete) //预约座位，查看预约历史

	// Web
	router.HandleFunc("/admin/login", controllers.AdminLogin).Methods(http.MethodPost,http.MethodOptions) //管理员登陆
	router.HandleFunc("/admin/token", controllers.VerifyToken).Methods(http.MethodGet) //token验证
	router.HandleFunc("/admin/statistic/init",controllers.InitHandler).Methods(http.MethodGet,http.MethodOptions) //获得统计数据
	router.HandleFunc("/admin/statistic/week", controllers.WeekHandler).Methods(http.MethodGet,http.MethodOptions) //获得一周预约数据
	router.HandleFunc("/admin/floor", controllers.FloorHandler).Methods(http.MethodGet,http.MethodPost,http.MethodOptions) //获得所有楼层，创建楼层
	router.HandleFunc("/admin/floor/{id}", controllers.FloorDelete).Methods(http.MethodDelete,http.MethodOptions)	//删除楼层
	router.HandleFunc("/admin/area", controllers.AreaGetByFloor).Methods(http.MethodGet,http.MethodOptions) //获取指定楼层区域
	router.HandleFunc("/admin/area", controllers.AreaCreate).Methods(http.MethodPost,http.MethodOptions)  //创建区域
	router.HandleFunc("/admin/area/update", controllers.AreaUpdate).Methods(http.MethodPost,http.MethodOptions)	//修改区域信息
	router.HandleFunc("/admin/area/{id}", controllers.AreaDelete).Methods(http.MethodDelete,http.MethodOptions)	//　删除区域
	router.HandleFunc("/admin/seat", controllers.SeatGetByArea).Methods(http.MethodGet,http.MethodOptions) //获取指定区域座位
	router.HandleFunc("/admin/seat", controllers.SeatCreate).Methods(http.MethodPost,http.MethodOptions)  //创建座位
	router.HandleFunc("/admin/seat/enable/{id}", controllers.SeatEnable).Methods(http.MethodGet,http.MethodOptions)  //启用座位
	router.HandleFunc("/admin/seat/disable/{id}", controllers.SeatDisable).Methods(http.MethodGet,http.MethodOptions)  //禁用座位
	router.HandleFunc("/admin/seat/update", controllers.SeatUpdate).Methods(http.MethodPost,http.MethodOptions)	//修改座位信息
	router.HandleFunc("/admin/seat/{id}", controllers.SeatDelete).Methods(http.MethodDelete,http.MethodOptions)	//删除座位
	router.HandleFunc("/admin/student", controllers.StudentList).Methods(http.MethodGet,http.MethodOptions) //获得学生列表
	router.HandleFunc("/admin/student/enable/{id}", controllers.StudentEnable).Methods(http.MethodGet,http.MethodOptions)  //用户正常
	router.HandleFunc("/admin/student/disable/{id}", controllers.StudentDisable).Methods(http.MethodGet,http.MethodOptions)  //用户封禁
	router.HandleFunc("/admin/student", controllers.StudentCreate).Methods(http.MethodPost,http.MethodOptions) //创建学生
	router.HandleFunc("/admin/student/update", controllers.StudentUpdate).Methods(http.MethodPost,http.MethodOptions)	//修改学生信息
	router.HandleFunc("/admin/student/{id}", controllers.StudentDelete).Methods(http.MethodDelete,http.MethodOptions)	//删除学生
	router.HandleFunc("/admin/student/password/{id}", controllers.StudentPWD).Methods(http.MethodGet,http.MethodOptions)	//重制密码
	router.HandleFunc("/admin/notice", controllers.NoticesHandler).Methods(http.MethodPost, http.MethodGet,http.MethodOptions)  //创建通知，获得通知列表
	router.HandleFunc("/admin/notice/{id}", controllers.NoticeHandler).Methods(http.MethodGet,http.MethodDelete,http.MethodOptions)	//获得指定通知，删除通知
	router.HandleFunc("/admin/notice/update", controllers.NoticeUpdate).Methods(http.MethodPost,http.MethodOptions)	//修改通知
	router.HandleFunc("/admin/appointment", controllers.AppointmentList).Methods(http.MethodGet,http.MethodOptions)	//获取预约


	router.HandleFunc("/floor", controllers.FloorHandler).Methods(http.MethodPost, http.MethodGet)  //创建楼层，获得所有楼层
	router.HandleFunc("/floor/{id}", controllers.FloorDelete).Methods(http.MethodDelete)	//删除楼层
	//router.HandleFunc("/floor/update", controllers.FloorUpdate).Methods(http.MethodPost)	//修改楼层信息


	router.HandleFunc("/area/{floor_id}", controllers.AreaGetByFloor).Methods(http.MethodGet) //获取指定楼层区域




	router.HandleFunc("/seat/{area_id}", controllers.SeatGetByArea).Methods(http.MethodGet) //获取指定区域座位



	router.HandleFunc("/notices", controllers.NoticesHandler).Methods(http.MethodPost, http.MethodGet)  //创建通知，获得通知列表

	router.HandleFunc("/notice/find_by_name", controllers.NoticeFindByName).Queries("name", "{name}").Methods(http.MethodGet)	//修改通知
	router.HandleFunc("/notice/find_by_username", controllers.NoticeFindByUsername).Queries("username", "{username}").Methods(http.MethodGet)	//修改通知

	router.HandleFunc("/student/create", controllers.StudentCreate).Queries("state", "{state}").Methods(http.MethodPost) //创建学生

	router.HandleFunc("/student/{id}", controllers.StudentHandler).Methods(http.MethodGet,http.MethodDelete) //创建学生,删除学生

	router.StrictSlash(true)
	router.Use(middlewares.Auth)
	conf := utils.GetConf()
	port := conf.Basic.Port

	var (
		err error
	)
	if conf.Basic.Tlsconf == nil {
		msg := fmt.Sprintf("Serving device manager on http://0.0.0.0:%d", port)
		utils.Log.Infoln(msg)
		fmt.Println(msg)
		err = http.ListenAndServe(fmt.Sprintf( ":%d", port), router)
	} else {
		msg := fmt.Sprintf("Serving device manager on https://0.0.0.0:%d", port)
		utils.Log.Infoln(msg)
		fmt.Println(msg)
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", port), conf.Basic.Tlsconf.Certfile, conf.Basic.Tlsconf.Keyfile, router)
	}

	if err != nil {
		panic(err)
	}
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	w.Write([]byte("Cors Request"))
}


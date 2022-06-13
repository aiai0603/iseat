package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"lib/models"
	"lib/utils"
	"net/http"
	"net/url"
	"strconv"
)
func StudentAppointmentsHandle(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		id,_:= strconv.Atoi(r.Header.Get("id"))
		stu , err:= models.GetStudent(uint(id))
		if err != nil{
			msg := message.GetMsg(utils.MsgAppointment, "get_student")
			http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
			return
		}
		if stu.State == 0{
			msg := message.GetMsg(utils.MsgAppointment, "student_banned")
			http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
			return
		}
		app := models.StuAppointment{}
		app.StudentID = stu.ID
		err = json.NewDecoder(r.Body).Decode(&app)
		if err != nil {
			msg := message.GetMsg(utils.MsgAppointment, "invalid_appointment")
			http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
			return
		}
		err = app.StudentCreateAppointment(stu,message)
		if err != nil{
			msg := message.GetMsg(utils.MsgAppointment, "create_appointment")
			http.Error(w, fmt.Sprintf(msg,err), http.StatusInternalServerError)
			return
		}
		utils.JsonResponse(app, w, http.StatusOK)
	case http.MethodGet:
		id,_:= strconv.Atoi(r.Header.Get("id"))
		stu , err:= models.GetStudent(uint(id))
		if err != nil{
			msg := message.GetMsg(utils.MsgAppointment, "get_student")
			http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
			return
		}
		apps,err:= models.StudentGetAppointmentList(stu)
		if err != nil{
			msg := message.GetMsg(utils.MsgAppointment, "get_appointment_list")
			http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
			return
		}
		utils.JsonResponse(apps, w, http.StatusOK)
	}
}

func StudentAppointmentIn(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	id,_:= strconv.Atoi(r.Header.Get("id"))
	stu , err:= models.GetStudent(uint(id))
	if err != nil{
		msg := message.GetMsg(utils.MsgAppointment, "get_student")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	apps,err:= models.StudentGetAppointmentInList(stu)
	if err != nil{
		msg := message.GetMsg(utils.MsgAppointment, "get_appointment_list")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(apps, w, http.StatusOK)
}

func StudentDelAppointment(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgAppointment, "invalid_appointment")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.DeleteAppointment(uint(appID))
	if err != nil {
		msg := message.GetMsg(utils.MsgAppointment, "delete_appointment")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse("ok",w, http.StatusOK)
}

func CreateAppointment(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json


	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	token := r.Header.Get("token")
	if token == ""{
		msg := message.GetMsg(utils.MsgToken, "no_token")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}
	claims, err := utils.ParseToken(token)
	if err != nil{
		msg := message.GetMsg(utils.MsgToken, "token")
		http.Error(w, fmt.Sprintf(msg,err), http.StatusUnauthorized)
		return
	}
	if claims.(jwt.MapClaims)["admin"] != false{
		msg := message.GetMsg(utils.MsgToken, "token_student")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}

	app := &models.Appointment{StudentID:uint(claims.(jwt.MapClaims)["id"].(float64))}
	err = json.NewDecoder(r.Body).Decode(app)
	if err != nil {
		msg := message.GetMsg(utils.MsgAppointment, "invalid_appointment")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	err = app.CreateAppointment()
	if err != nil {
		msg := message.GetMsg(utils.MsgAppointment, "create_appointment")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(app, w, http.StatusCreated)
}

func AppointmentList(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()
	s, _ :=url.QueryUnescape(r.URL.RawQuery)
	filer := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &filer)

	data,count,err := models.GetAppointmentList(filer)
	if err != nil{
		msg := message.GetMsg(utils.MsgAppointment, "invalid_appointment")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	result := models.Table{Code: 0,Count: count,Data: data}
	utils.JsonResponse(result,w, http.StatusOK)
}

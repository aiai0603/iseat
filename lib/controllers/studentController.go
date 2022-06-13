package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"lib/models"
	"lib/utils"
	"net/http"
	"net/url"
	"strconv"
)

func StudentLogin(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	studentLogin := &models.Login{}
	err := json.NewDecoder(r.Body).Decode(studentLogin)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_student")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	if studentLogin.Username == "" || studentLogin.Password == ""{
		msg := message.GetMsg(utils.MsgStudent, "invalid_student")
		http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
		return
	}

	student,err := studentLogin.GetStudent()
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "get_student")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	if student.Name == ""{
		msg := message.GetMsg(utils.MsgStudent, "student_username")
		http.Error(w, fmt.Sprintf(msg), http.StatusBadRequest)
		return
	}
	if student.Password != utils.Md5V2(studentLogin.Password) {
		msg := message.GetMsg(utils.MsgStudent, "student_password")
		http.Error(w, fmt.Sprintf(msg), http.StatusBadRequest)
		return
	}

	token, err := utils.CreateToken(student.ID, false)
	if err != nil{
		msg := message.GetMsg(utils.MsgToken, "create_token")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	student.Password = ""
	res := models.ResStudentLogin{Student: student, Token: token}
	utils.JsonResponse(res, w, http.StatusOK)
}

func StudentInf(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	id,_:= strconv.Atoi(r.Header.Get("id"))
	stu , err:= models.GetStudent(uint(id))
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "get_student")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	count,err := stu.GetAppointmentCount()
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "get_appointment_count")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	stu.Password = ""
	result := models.StudentInf{
		Student:stu,
		Count: count,
	}
	utils.JsonResponse(result, w, http.StatusOK)
}

func StudentChangePwd(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	id,_:= strconv.Atoi(r.Header.Get("id"))
	stu , err:= models.GetStudent(uint(id))
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "get_student")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}

	pwd := models.StudentPwd{}
	err = json.NewDecoder(r.Body).Decode(&pwd)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_pwd")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	if stu.Password != utils.Md5V2(pwd.OldPwd){
		msg := message.GetMsg(utils.MsgStudent, "error_old_password")
		http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
		return
	}

	err = stu.ChangePwd(pwd.NewPwd)
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "change_password")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	stu.Password = ""
	utils.JsonResponse(stu, w, http.StatusOK)
}

func StudentHandler(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	switch r.Method {
	case http.MethodDelete:
		err := models.DeleteStudent(uint(id))
		if err == nil{
			msg := message.GetMsg(utils.MsgStudent, "delete_student")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodGet:
		stu, err := models.GetStudent(uint(id))
		if err == nil{
			msg := message.GetMsg(utils.MsgStudent, "get_student")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		utils.JsonResponse(stu, w, http.StatusOK)
	}
}

func StudentCreate(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	student := &models.Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_student")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	err = student.CreateStudent()
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "create_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(student, w, http.StatusCreated)


}

func StudentList(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	result := models.Table{}
	filer := make(map[string]interface{})
	s, _ :=url.QueryUnescape(r.URL.RawQuery)
	err := json.Unmarshal([]byte(s), &filer)

	students,count,err := models.GetStudentList(filer)
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "page_parameters")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	result.Code = 0
	result.Data = students
	result.Count = count
	utils.JsonResponse(result, w, http.StatusOK)
}

func StudentEnable(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.StudentState(uint(id),1)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}

func StudentDisable(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.StudentState(uint(id),0)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}

func StudentUpdate(w http.ResponseWriter, r *http.Request) {
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	student := &models.Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	err = student.UpdateStudent()
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "update_student")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(student, w, http.StatusOK)
}

func StudentDelete(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.DeleteStudent(uint(id))
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "delete_seat")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}

func StudentPWD(w http.ResponseWriter, r *http.Request) {
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgStudent, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.UpdateStudentPWD(uint(id))
	if err != nil{
		msg := message.GetMsg(utils.MsgStudent, "update_student")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(id, w, http.StatusOK)
}
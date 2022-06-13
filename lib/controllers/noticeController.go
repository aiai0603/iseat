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

func StudentGetNotices(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	list, err := models.StudentGetNoticeList()
	if err != nil {
		msg := message.GetMsg(utils.MsgNotice, "student_get_notices")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}

	utils.JsonResponse(list, w, http.StatusOK)

}



func NoticesHandler(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		id,_:= strconv.Atoi(r.Header.Get("id"))
		notice := models.Notice{}
		err := json.NewDecoder(r.Body).Decode(&notice)
		if err != nil {
			msg := message.GetMsg(utils.MsgNotice, "invalid_notice")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		notice.AdminID=uint(id)
		err = notice.CreateNotice()
		if err != nil {
			msg := message.GetMsg(utils.MsgNotice, "create_notice")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		utils.JsonResponse(notice, w, http.StatusOK)
	case http.MethodGet:
		result := models.Table{}
		filer := make(map[string]interface{})
		s, _ :=url.QueryUnescape(r.URL.RawQuery)
		err := json.Unmarshal([]byte(s), &filer)

		notices,count ,err := models.ListNotices(filer)
		if err != nil {
			msg := message.GetMsg(utils.MsgNotice, "get_list")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
			return
		}
		result.Code = 0
		result.Data = notices
		result.Count = count
		utils.Respond(result, err, w)
	}
}

func NoticeHandler(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgNotice, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		notice,err := models.GetNotice(uint(id))
		if err != nil{
			msg := message.GetMsg(utils.MsgNotice, "get_notice")
			http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
			return
		}
		err = models.DeleteNotice(notice)
		if err != nil{
			msg := message.GetMsg(utils.MsgNotice, "delete_notice")
			http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		notice, err := models.GetNotice(uint(id))
		if err != nil {
			msg := message.GetMsg(utils.MsgNotice, "get_notice")
			http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
			return
		}
		utils.JsonResponse(notice, w, http.StatusOK)
	}
}

func NoticeUpdate(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	notice := models.Notice{}
	err := json.NewDecoder(r.Body).Decode(&notice)
	if err != nil {
		msg := message.GetMsg(utils.MsgAdmin, "invalid_admin")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	err = notice.UpdateNotice()
	if err != nil{
		msg := message.GetMsg(utils.MsgAdmin, "update_notice")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(notice, w, http.StatusOK)
}

func NoticeFindByName(w http.ResponseWriter, r *http.Request){
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
	if claims.(jwt.MapClaims)["admin"] != true{
		msg := message.GetMsg(utils.MsgToken, "token_admin")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}


}

func NoticeFindByUsername(w http.ResponseWriter, r *http.Request){
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
	if claims.(jwt.MapClaims)["admin"] != true{
		msg := message.GetMsg(utils.MsgToken, "token_admin")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}


}


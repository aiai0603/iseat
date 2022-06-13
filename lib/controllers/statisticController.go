package controllers

import (
	"fmt"
	"lib/models"
	"lib/utils"
	"net/http"
)

func InitHandler(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	result,err :=models.GetCount()
	if err != nil {
		msg := message.GetMsg(utils.MsgStatistic, "get_init")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(result, w, http.StatusOK)
}

func WeekHandler(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	result,err :=models.GetWeek()
	if err != nil{
		msg := message.GetMsg(utils.MsgStatistic, "get_init")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(result, w, http.StatusOK)
}
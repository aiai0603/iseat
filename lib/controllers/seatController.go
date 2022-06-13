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

func SeatCreate(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	seat := &models.Seat{}
	err := json.NewDecoder(r.Body).Decode(seat)
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	err = seat.CreateSeat()
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "create_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(seat, w, http.StatusOK)
}

func SeatGetByArea(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()
	result := models.Table{}
	filer := make(map[string]interface{})
	s, _ :=url.QueryUnescape(r.URL.RawQuery)
	err := json.Unmarshal([]byte(s), &filer)
	if filer["area_id"] == nil{
		result.Code = 1
		result.Msg = "请先选择楼层和区域"
		utils.JsonResponse(result,w,http.StatusOK)
		return
	}
	seats,count,err := models.GetSeatByArea(filer)
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "get_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	result.Count = count
	result.Code = 0
	result.Data = seats
	utils.JsonResponse(result, w, http.StatusOK)
}

func SeatUpdate(w http.ResponseWriter, r *http.Request) {
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	seat := &models.Seat{}
	err := json.NewDecoder(r.Body).Decode(seat)
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	err = seat.UpdateSeat()
	if err != nil{
		msg := message.GetMsg(utils.MsgSeat, "update_seat")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(seat, w, http.StatusOK)
}

func SeatEnable(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.SeatState(uint(id),1)
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}

func SeatDisable(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.SeatState(uint(id),0)
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}

func SeatDelete(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.DeleteSeat(uint(id))
	if err != nil {
		msg := message.GetMsg(utils.MsgSeat, "delete_seat")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(id,w,http.StatusOK)
}
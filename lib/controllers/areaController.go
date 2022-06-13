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

func  StudentGetArea(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	app := models.Appointment{}
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "invalid_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	floorID, err := strconv.Atoi(vars["floor_id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "invalid_floor_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	areas, err := models.StudentGetAreaByFloor(uint(floorID))
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "get_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	var result []models.StudentArea
	for _,area := range areas{
		stuArea :=models.StudentArea{Area:area}
		err := stuArea.GetAreaInf(app)
		if err != nil{
			msg := message.GetMsg(utils.MsgArea, "get_area_inf")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
			return
		}
		result = append(result,stuArea)
	}
	utils.JsonResponse(result, w, http.StatusOK)
}

func AreaCreate(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	area := &models.Area{}
	err := json.NewDecoder(r.Body).Decode(area)
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "invalid_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	err = area.CreateArea()
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "create_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(area, w, http.StatusOK)
}

func AreaGetByFloor(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()
	result := models.Table{}
	filer := make(map[string]interface{})
	s, _ :=url.QueryUnescape(r.URL.RawQuery)
	err := json.Unmarshal([]byte(s), &filer)
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "get_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	if filer["searchParams"] == nil {
		result.Code = 1
		result.Msg="请先选择楼层"
		result.Data=[]models.Area{}
	}else{
		areas, err := models.GetAreaByFloor(filer)
		if err != nil {
			msg := message.GetMsg(utils.MsgArea, "get_area")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		result.Code = 0
		result.Data=areas
	}


	utils.JsonResponse(result, w, http.StatusOK)
}

func AreaUpdate(w http.ResponseWriter, r *http.Request) {
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	area := &models.Area{}
	err := json.NewDecoder(r.Body).Decode(area)
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "invalid_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	err = area.UpdateArea()
	if err != nil{
		msg := message.GetMsg(utils.MsgArea, "update_area")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(area, w, http.StatusOK)
}

func AreaDelete(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	err = models.DeleteArea(uint(id))
	if err != nil {
		msg := message.GetMsg(utils.MsgArea, "delete_area")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(id, w, http.StatusOK)
}
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

func FloorHandler(w http.ResponseWriter, r *http.Request)  {
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()


	switch r.Method {
	case http.MethodPost:
		f := &models.Floor{}
		err := json.NewDecoder(r.Body).Decode(f)
		if err != nil {
			msg := message.GetMsg(utils.MsgFloor, "invalid_floor")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		err = f.CreateFloor()
		if err != nil {
			msg := message.GetMsg(utils.MsgFloor, "create_floor")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		utils.JsonResponse(f, w, http.StatusOK)
	case http.MethodGet:
		filer := make(map[string]interface{})
		s, _ :=url.QueryUnescape(r.URL.RawQuery)
		err := json.Unmarshal([]byte(s), &filer)
		floors,err := models.GetAllFloors(filer)
		if err != nil{
			msg := message.GetMsg(utils.MsgFloor, "get_floor")
			http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
			return
		}
		result := models.Table{Data: floors,Code: 0}
		utils.JsonResponse(result, w, http.StatusOK)
	}
}

func FloorDelete(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := message.GetMsg(utils.MsgFloor, "invalid_id")
		http.Error(w, fmt.Sprintf(msg, err), http.StatusBadRequest)
		return
	}
	floor, err := models.GetFloors(uint(id), message)
	if err != nil {
		msg := message.GetMsg(utils.MsgFloor, "get_floor")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	err = models.DeleteFloor(floor)
	if err != nil {
		msg := message.GetMsg(utils.MsgFloor, "delete_floor")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func FloorUpdate(w http.ResponseWriter, r *http.Request){
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

	f := &models.Floor{}
	err = json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		msg := message.GetMsg(utils.MsgFloor, "invalid_floor")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}

	err = f.UpdateFloor()
	if err != nil{
		msg := message.GetMsg(utils.MsgFloor, "update_floor")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusBadRequest)
		return
	}
	utils.JsonResponse(f, w, http.StatusOK)
}



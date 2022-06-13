package controllers

import (
	"fmt"
	"lib/utils"
	"net/http"
)

func VerifyToken(w http.ResponseWriter, r *http.Request){
	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	token := r.Header.Get("token")
	if token == ""{
		msg := message.GetMsg(utils.MsgToken, "no_token")
		http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
		return
	}
	_, err := utils.ParseToken(token)
	if err != nil{
		msg := message.GetMsg(utils.MsgToken, "token")
		http.Error(w, fmt.Sprintf(msg,err), http.StatusUnauthorized)
		return
	}
	utils.JsonResponse(token, w, http.StatusOK)
}

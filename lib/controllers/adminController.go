package controllers

import (
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"lib/models"
	"lib/utils"
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json
	if r.Method == http.MethodOptions {
		return
	}

	message := new(utils.Message)
	message.SetLanguage(r)
	defer r.Body.Close()

	adminLogin := &models.Login{}
	err := json.NewDecoder(r.Body).Decode(adminLogin)
	if err != nil {
		msg := message.GetMsg(utils.MsgAdmin, "invalid_admin")
		http.Error(w, fmt.Sprintf(msg, err.Error()), http.StatusInternalServerError)
		return
	}
	if adminLogin.Username == "" || adminLogin.Password == ""{
		msg := message.GetMsg(utils.MsgAdmin, "invalid_admin_login")
		http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
		return
	}

	admin,err := adminLogin.GetAdmin()
	if err != nil{
		msg := message.GetMsg(utils.MsgAdmin, "get_admin")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}
	if admin.Name == ""{
		msg := message.GetMsg(utils.MsgAdmin, "admin_username")
		http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
		return
	}
	if admin.Password !=utils.Md5V2(adminLogin.Password)  {
		msg := message.GetMsg(utils.MsgAdmin, "admin_password")
		http.Error(w, fmt.Sprintf(msg), http.StatusInternalServerError)
		return
	}

	token, err := utils.CreateToken(admin.ID, true)
	if err != nil{
		msg := message.GetMsg(utils.MsgToken, "create_token")
		http.Error(w, fmt.Sprintf(msg,err.Error()), http.StatusInternalServerError)
		return
	}

	res := models.ResAdminLogin{Admin: admin, Token: token}

	utils.JsonResponse(res, w, http.StatusOK)
}

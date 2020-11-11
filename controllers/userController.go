package controllers

import (
	"encoding/json"
	"internship_project/services"
	"internship_project/utils"
	"net/http"
)

type UserController struct {
	Service services.UserService
}

func (controller *UserController) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("id_token")

	user, err := controller.Service.GoogleSignIn(token)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
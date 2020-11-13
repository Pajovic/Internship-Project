package controllers

import (
	"fmt"
	"internship_project/services"
	"internship_project/utils"
	"net/http"

	"github.com/markbates/goth/gothic"
)

type UserController struct {
	Service services.UserService
}

func (controller *UserController) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	u, err := controller.Service.GoogleSignIn(user)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	fmt.Println("\nYou are logged in as:")
	fmt.Println("     ", u)
}

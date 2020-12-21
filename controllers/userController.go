package controllers

import (
	"encoding/json"
	"errors"
	"internship_project/services"
	"internship_project/utils"
	"net/http"
)

type UserController struct {
	Service services.UserService
}

func (controller *UserController) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Token string
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		utils.WriteErrToClient(w, errors.New("Couldn't decode parameters"))
		return
	}

	u, err := controller.Service.GoogleSignIn(params.Token)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	jwt, err := utils.CreateJWT(u)
	/* fmt.Println("\nYou are logged in as:", u)
	fmt.Println("Your JWT is: ", jwt, "\n") */
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwt)
}

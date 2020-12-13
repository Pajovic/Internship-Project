package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"internship_project/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type ShopController struct {
	Service services.ShopService
}

func (controller *ShopController) GetAllShops(w http.ResponseWriter, r *http.Request) {
	shops, err := controller.Service.GetAllShops()
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
}

func (controller *ShopController) GetShopById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	shop, err := controller.Service.GetShop(idParam)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shop)
}

func (controller *ShopController) AddShop(w http.ResponseWriter, r *http.Request) {
	var newShop models.Shop
	json.NewDecoder(r.Body).Decode(&newShop)
	err := controller.Service.AddNewShop(&newShop)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newShop)
}

func (controller *ShopController) UpdateShop(w http.ResponseWriter, r *http.Request) {
	var updateShop models.Shop
	json.NewDecoder(r.Body).Decode(&updateShop)

	err := controller.Service.UpdateShop(updateShop)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateShop)
}

func (controller *ShopController) DeleteShop(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	err := controller.Service.DeleteShop(idParam)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}

func (controller *ShopController) GetAddress(w http.ResponseWriter, r *http.Request) {
	var shopId string = mux.Vars(r)["id"]

	address, err := controller.Service.GetAddress(shopId)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

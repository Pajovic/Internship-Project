package controllers

import (
	"encoding/json"
	"internship_project/errorhandler"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
)

type ProductController struct {
	Service services.ProductService
}

func (controller *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := controller.Service.GetAllProducts()
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (controller *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	product, err := controller.Service.GetProduct(idParam)
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (controller *ProductController) AddProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	json.NewDecoder(r.Body).Decode(&newProduct)
	err := controller.Service.AddNewProduct(&newProduct)
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

func (controller *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	var updateProduct models.Product
	json.NewDecoder(r.Body).Decode(&updateProduct)

	updateProduct.ID = idParam

	err := controller.Service.Updateproduct(updateProduct)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateProduct)

}

func (controller *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	err := controller.Service.DeleteProduct(idParam)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}

func writeErrToClient(w http.ResponseWriter, err error) {
	errMsg, code := errorhandler.GetErrorMsg(err)
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}

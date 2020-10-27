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
	idEmployee := r.Header.Get("employeeID")
	products, err := controller.Service.GetAllProducts(idEmployee)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (controller *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	idEmployee := r.Header.Get("employeeID")

	product, err := controller.Service.GetProduct(idParam, idEmployee)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (controller *ProductController) AddProduct(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.Header.Get("employeeID")

	var newProduct models.Product
	json.NewDecoder(r.Body).Decode(&newProduct)
	err := controller.Service.AddNewProduct(&newProduct, idEmployee)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

func (controller *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.Header.Get("employeeID")

	var updateProduct models.Product
	json.NewDecoder(r.Body).Decode(&updateProduct)

	err := controller.Service.UpdateProduct(updateProduct, idEmployee)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateProduct)

}

func (controller *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	idEmployee := r.Header.Get("employeeID")

	err := controller.Service.DeleteProduct(idParam, idEmployee)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}

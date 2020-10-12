package controllers

import (
	"encoding/json"
	"fmt"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
)

type ProductController struct {
	Service         services.ProductService
	EmployeeService services.EmployeeService
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
	idEmployee := mux.Vars(r)["employeeId"]

	product, err := controller.Service.GetProduct(idParam)

	employeeRights := controller.Service.GetEmployeePermissions(idEmployee, idParam)

	fmt.Println("Controller fmt")
	fmt.Println(employeeRights)

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

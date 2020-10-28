package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"internship_project/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//EmployeeController .
type EmployeeController struct {
	Service services.EmployeeService
}

// GetAllEmployees is used for getting all employees from the database
func (controller *EmployeeController) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.Header.Get("employeeID")
	allEmployees, err := controller.Service.GetAllEmployees(idEmployee)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(allEmployees)
}

// AddNewEmployee is used to add a new employee
func (controller *EmployeeController) AddNewEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&newEmployee)

	err := controller.Service.AddNewEmployee(&newEmployee)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func (controller *EmployeeController) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.Header.Get("employeeID")
	id := mux.Vars(r)["id"] // Because ID is string in database

	employee, err := controller.Service.GetEmployeeByID(id, idEmployee)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee is used to update employee's info
func (controller *EmployeeController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var updatedEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&updatedEmployee)

	err := controller.Service.UpdateEmployee(updatedEmployee)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}

// DeleteEmployee is used to delete employee
func (controller *EmployeeController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := controller.Service.DeleteEmployee(id)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}

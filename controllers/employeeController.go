package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllEmployees is used for getting all employees from the database
func GetAllEmployees(w http.ResponseWriter, r *http.Request, e []models.Employee) {
	allEmployees := services.GetAllEmployees(e)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(allEmployees)
}

// AddNewEmployee is used to add a new employee
func AddNewEmployee(w http.ResponseWriter, r *http.Request, e *[]models.Employee) {
	var newEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&newEmployee)

	e, err := services.AddNewEmployee(e, newEmployee)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("An error has occured while adding employee."))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(e)
}

// GetEmployeeByID is used to find a specific employee
func GetEmployeeByID(w http.ResponseWriter, r *http.Request, e []models.Employee) {
	idString := mux.Vars(r)["id"]   // Because ID is string in database
	id, _ := strconv.Atoi(idString) // Temporarily, for now, to check if ID is invalid (larger than slice length)

	if id < 0 || id > len(e) {
		w.WriteHeader(404)
		w.Write([]byte("Invalid Employee ID"))
		return
	}

	employee := services.GetEmployeeByID(e, idString)

	if employee == (models.Employee{}) {
		w.WriteHeader(404)
		w.Write([]byte("There is no Employee with ID " + idString))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee is used to update employee's info
func UpdateEmployee(w http.ResponseWriter, r *http.Request, e *[]models.Employee) {
	var updatedEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&updatedEmployee)

	idString := updatedEmployee.ID
	id, _ := strconv.Atoi(idString)

	if id < 0 || id > len(*e) {
		w.WriteHeader(404)
		w.Write([]byte("Invalid Employee ID"))
		return
	}

	*e = services.UpdateEmployee(e, updatedEmployee, idString)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(e)
}

// DeleteEmployee is used to delete employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request, e *[]models.Employee) {
	idString := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idString)

	if id < 0 || id > len(*e) {
		w.WriteHeader(404)
		w.Write([]byte("Invalid Employee ID"))
		return
	}

	*e = services.DeleteEmployee(*e, idString)

	w.WriteHeader(204)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(*e)
}

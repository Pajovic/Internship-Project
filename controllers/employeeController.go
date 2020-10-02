package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// GetAllEmployees is used for getting all employees from the database
func GetAllEmployees(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	allEmployees, err := services.GetAllEmployees(conn)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("An error has occured while getting all employees."))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(allEmployees)
}

// AddNewEmployee is used to add a new employee
func AddNewEmployee(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	var newEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&newEmployee)

	err := services.AddNewEmployee(conn, &newEmployee)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("An error has occured while adding employee."))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func GetEmployeeByID(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	id := mux.Vars(r)["id"] // Because ID is string in database

	employee, err := services.GetEmployeeByID(conn, id)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("An error has occured while getting the employee"))
		return
	}

	if employee == (models.Employee{}) {
		w.WriteHeader(404)
		w.Write([]byte("There is no Employee with ID " + id))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee is used to update employee's info
func UpdateEmployee(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	var updatedEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&updatedEmployee)

	err := services.UpdateEmployee(conn, updatedEmployee)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("An error has occured while updating."))
		return
	}
	w.WriteHeader(200)
}

// DeleteEmployee is used to delete employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	id := mux.Vars(r)["id"]

	err := services.DeleteEmployee(conn, id)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("An error has occured while deleting."))
		return
	}
	w.WriteHeader(204)
}

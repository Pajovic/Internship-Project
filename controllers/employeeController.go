package controllers

import (
	"encoding/json"
	"internship_project/errorhandler"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

// GetAllEmployees is used for getting all employees from the database
func GetAllEmployees(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	allEmployees, err := services.GetAllEmployees(conn)

	if err != nil {
		writeErrToClient(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(allEmployees)
}

// AddNewEmployee is used to add a new employee
func AddNewEmployee(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var newEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&newEmployee)

	err := services.AddNewEmployee(conn, &newEmployee)

	if err != nil {
		writeErrToClient(w, err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func GetEmployeeByID(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	id := mux.Vars(r)["id"] // Because ID is string in database

	employee, err := services.GetEmployeeByID(conn, id)

	if err != nil {
		writeErrToClient(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee is used to update employee's info
func UpdateEmployee(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var updatedEmployee models.Employee
	json.NewDecoder(r.Body).Decode(&updatedEmployee)

	err := services.UpdateEmployee(conn, updatedEmployee)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}

// DeleteEmployee is used to delete employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	id := mux.Vars(r)["id"]

	err := services.DeleteEmployee(conn, id)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}

func writeErrToClient(w http.ResponseWriter, err error) {
	errMsg, code := errorhandler.GetErrorMsg(err)
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}

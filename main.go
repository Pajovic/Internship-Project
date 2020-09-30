package main

import (
	"fmt"
	"internship_project/controllers"
	"internship_project/models"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var employees []models.Employee = []models.Employee{}
	r := mux.NewRouter()

	employeeRouter := r.PathPrefix("/employees").Subrouter()

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllEmployees(w, r, employees)
	}).Methods("GET")

	employeeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetEmployeeByID(w, r, employees)
	}).Methods("GET")

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddNewEmployee(w, r, &employees)
	}).Methods("POST")

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateEmployee(w, r, &employees)
	}).Methods("PUT")

	employeeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteEmployee(w, r, &employees)
	}).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world from Gorilla Mux")
}

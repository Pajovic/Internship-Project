package main

import (
	"Internship-Project/controllers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/company", controllers.GetAllCompanies).Methods("GET")

	r.HandleFunc("/company", controllers.AddCompany).Methods("POST")

	r.HandleFunc("/company/{id}", controllers.UpdateCompany).Methods("PUT")

	r.HandleFunc("/company/{id}", controllers.DeleteCompany).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world from Gorilla Mux")
}

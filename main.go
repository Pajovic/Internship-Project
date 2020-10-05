package main

import (
	"context"
	"fmt"
	"internship_project/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
)

type config struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DatabaseURL string `json:"database_url"`
}

func main() {
	var conf config
	if _, err := confl.DecodeFile("database.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.DatabaseURL)

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	r := mux.NewRouter()

	// Employee Routes
	employeeRouter := r.PathPrefix("/employees").Subrouter()

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllEmployees(w, r, db)
	}).Methods("GET")

	employeeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetEmployeeByID(w, r, db)
	}).Methods("GET")

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddNewEmployee(w, r, db)
	}).Methods("POST")

	employeeRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateEmployee(w, r, db)
	}).Methods("PUT")

	employeeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteEmployee(w, r, db)
	}).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

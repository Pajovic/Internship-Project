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

	connection, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	r := mux.NewRouter()

	companyRouter := r.PathPrefix("/company").Subrouter()

	companyRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllCompanies(w, r, connection)
	}).Methods("GET")

	companyRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetCompanyById(w, r, connection)
	}).Methods("GET")

	companyRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddCompany(w, r, connection)
	}).Methods("POST")

	companyRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateCompany(w, r, connection)
	}).Methods("PUT")

	companyRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteCompany(w, r, connection)
	}).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

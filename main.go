package main

import (
	"context"
	"fmt"
	"internship_project/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
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

	connection, err := pgx.Connect(context.Background(), conf.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close(context.Background())

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllCompanies(w, r, connection)
	}).Methods("GET")

	r.HandleFunc("/company/{id}", controllers.GetCompanyById).Methods("GET")

	r.HandleFunc("/company", controllers.AddCompany).Methods("POST")

	r.HandleFunc("/company/{id}", controllers.UpdateCompany).Methods("PUT")

	r.HandleFunc("/company/{id}", controllers.DeleteCompany).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world from Gorilla Mux")
}

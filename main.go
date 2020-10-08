package main

import (
	"context"
	"fmt"
	"internship_project/controllers"
	"internship_project/repositories"
	"internship_project/services"
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
	connpool := getConnectionPool()
	companyController := getController(connpool)
	employeeController := getEmployeeController(connpool)
	defer connpool.Close()

	r := mux.NewRouter()

	companyRouter := r.PathPrefix("/company").Subrouter()

	companyRouter.HandleFunc("", companyController.GetAllCompanies).Methods("GET")

	companyRouter.HandleFunc("/{id}", companyController.GetCompanyById).Methods("GET")

	companyRouter.HandleFunc("", companyController.AddCompany).Methods("POST")

	companyRouter.HandleFunc("/{id}", companyController.UpdateCompany).Methods("PUT")

	companyRouter.HandleFunc("/{id}", companyController.DeleteCompany).Methods("DELETE")

	// Employee Routes
	employeeRouter := r.PathPrefix("/employees").Subrouter()

	employeeRouter.HandleFunc("/", employeeController.GetAllEmployees).Methods("GET")
	employeeRouter.HandleFunc("/{id}", employeeController.GetEmployeeByID).Methods("GET")
	employeeRouter.HandleFunc("/", employeeController.AddNewEmployee).Methods("POST")
	employeeRouter.HandleFunc("/", employeeController.UpdateEmployee).Methods("PUT")
	employeeRouter.HandleFunc("/{id}", employeeController.DeleteEmployee).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

func getConnectionPool() *pgxpool.Pool {
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

	fmt.Println("Connected to database.")

	return connection
}

func getController(connpool *pgxpool.Pool) controllers.CompanyController {
	companyRepository := repositories.CompanyRepository{DB: connpool}
	companyService := services.CompanyService{Repository: companyRepository}
	companyController := controllers.CompanyController{Service: companyService}

	fmt.Println("Employee controller up and running.")

	return companyController
}

func getEmployeeController(connpool *pgxpool.Pool) controllers.EmployeeController {
	employeeRepository := repositories.EmployeeRepository{DB: connpool}
	employeeService := services.EmployeeService{Repository: employeeRepository}
	employeeController := controllers.EmployeeController{Service: employeeService}

	fmt.Println("Employee controller up and running.")

	return employeeController
}

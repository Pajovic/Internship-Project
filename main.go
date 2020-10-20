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
	employeeController := getEmployeeController(connpool)
	productController := getProductController(connpool, &employeeController.Service.Repository)
	companyController := getCompanyController(connpool)
	earController := getEarController(connpool)
	constraintController := getConstraintController(connpool)

	defer connpool.Close()

	r := mux.NewRouter()

	// Product Routes
	productRouter := r.PathPrefix("/product").Subrouter()
	productRouter.Headers("employeeID")

	productRouter.HandleFunc("", productController.GetAllProducts).Methods("GET")
	productRouter.HandleFunc("/{id}", productController.GetProductById).Methods("GET")
	productRouter.HandleFunc("", productController.AddProduct).Methods("POST")
	productRouter.HandleFunc("/{id}", productController.UpdateProduct).Methods("PUT")
	productRouter.HandleFunc("/{id}", productController.DeleteProduct).Methods("DELETE")

	// Company Routes
	companyRouter := r.PathPrefix("/company").Subrouter()
	companyRouter.Headers("companyID")

	companyRouter.HandleFunc("", companyController.GetAllCompanies).Methods("GET")
	companyRouter.HandleFunc("/{id}", companyController.GetCompanyById).Methods("GET")
	companyRouter.HandleFunc("", companyController.AddCompany).Methods("POST")
	companyRouter.HandleFunc("/{id}", companyController.UpdateCompany).Methods("PUT")
	companyRouter.HandleFunc("/{id}", companyController.DeleteCompany).Methods("DELETE")
	companyRouter.HandleFunc("/approve/{idear}", companyController.ApproveExternalAccess).Methods("PATCH")
	companyRouter.HandleFunc("/disapprove/{idear}", companyController.DisapproveExternalAccess).Methods("PATCH")

	// Employee Routes
	employeeRouter := r.PathPrefix("/employees").Subrouter()
	productRouter.Headers("employeeID")

	employeeRouter.HandleFunc("", employeeController.GetAllEmployees).Methods("GET")
	employeeRouter.HandleFunc("/{id}", employeeController.GetEmployeeByID).Methods("GET")
	employeeRouter.HandleFunc("", employeeController.AddNewEmployee).Methods("POST")
	employeeRouter.HandleFunc("", employeeController.UpdateEmployee).Methods("PUT")
	employeeRouter.HandleFunc("/{id}", employeeController.DeleteEmployee).Methods("DELETE")

	// External Access Rules Routes
	earRouter := r.PathPrefix("/ear").Subrouter()

	earRouter.HandleFunc("", earController.GetAllEars).Methods("GET")
	earRouter.HandleFunc("/{id}", earController.GetEarById).Methods("GET")
	earRouter.HandleFunc("", earController.AddEar).Methods("POST")
	earRouter.HandleFunc("", earController.UpdateEar).Methods("PUT")
	earRouter.HandleFunc("/{id}", earController.DeleteEar).Methods("DELETE")

	// Constraints Routes
	constraintRouter := r.PathPrefix("/constraint").Subrouter()

	constraintRouter.HandleFunc("", constraintController.GetAllConstraints).Methods("GET")
	constraintRouter.HandleFunc("/{id}", constraintController.GetConstraintById).Methods("GET")
	constraintRouter.HandleFunc("", constraintController.AddConstraint).Methods("POST")
	constraintRouter.HandleFunc("", constraintController.UpdateConstraint).Methods("PUT")
	constraintRouter.HandleFunc("/{id}", constraintController.DeleteConstraint).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

func getConnectionPool() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("dbconfig.conf", &conf); err != nil {
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

func getProductController(connpool *pgxpool.Pool, employeeRepo *repositories.EmployeeRepository) controllers.ProductController {

	productRepository := repositories.ProductRepository{DB: connpool}
	productService := services.ProductService{ProductRepository: productRepository, EmployeeRepository: *employeeRepo}
	productController := controllers.ProductController{Service: productService}

	fmt.Println("Product controller up and running.")

	return productController
}

func getCompanyController(connpool *pgxpool.Pool) controllers.CompanyController {
	companyRepository := repositories.CompanyRepository{DB: connpool}
	companyService := services.CompanyService{Repository: companyRepository}
	companyController := controllers.CompanyController{Service: companyService}

	fmt.Println("Company controller up and running.")

	return companyController
}

func getEmployeeController(connpool *pgxpool.Pool) controllers.EmployeeController {
	employeeRepository := repositories.EmployeeRepository{DB: connpool}
	employeeService := services.EmployeeService{Repository: employeeRepository}
	employeeController := controllers.EmployeeController{Service: employeeService}

	fmt.Println("\nEmployee controller up and running.")

	return employeeController
}

func getEarController(connpool *pgxpool.Pool) controllers.EarController {
	earRepository := repositories.EarRepository{DB: connpool}
	earService := services.EarService{Repository: earRepository}
	earController := controllers.EarController{Service: earService}

	fmt.Println("External access rights controller up and running.")

	return earController
}

func getConstraintController(connpool *pgxpool.Pool) controllers.ConstraintController {
	constraintRepository := repositories.ConstraintRepository{DB: connpool}
	constraintService := services.ConstraintService{Repository: constraintRepository}
	constraintController := controllers.ConstraintController{Service: constraintService}

	fmt.Println("Constraints controller up and running.")

	return constraintController
}

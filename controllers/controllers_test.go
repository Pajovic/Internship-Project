package controllers

import (
	"context"
	"fmt"
	"internship_project/repositories"
	"internship_project/services"
	"internship_project/utils"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
)

var (
	connpool *pgxpool.Pool

	CompanyCont       CompanyController
	ProductCont       ProductController
	EmployeeCont      EmployeeController
	ConstraintCont    ConstraintController
	ExternalRightCont ExternalRightController
)

type Config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

func TestMain(m *testing.M) {
	connpool = GetTestConnectionPool()
	defer connpool.Close()

	CompanyCont = GetCompanyController(connpool)
	EmployeeCont = GetEmployeeController(connpool)
	ProductCont = getProductController(connpool, &EmployeeCont.Service.Repository)
	ConstraintCont = GetConstraintController(connpool)
	ExternalRightCont = GetExternalRightController(connpool)

	utils.SetUpTables(connpool)

	os.Exit(m.Run())
}

func GetTestConnectionPool() *pgxpool.Pool {
	var conf Config
	if _, err := confl.DecodeFile("./../dbconfig.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.TestDatabaseURL)

	connection, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to test database.")

	return connection
}

func GetCompanyController(connpool *pgxpool.Pool) CompanyController {
	companyRepository := repositories.NewCompanyRepo(connpool)
	companyService := services.CompanyService{Repository: companyRepository}
	companyController := CompanyController{Service: companyService}

	fmt.Println("Company controller up and running.")

	return companyController
}

func GetConstraintController(connpool *pgxpool.Pool) ConstraintController {
	constraintRepository := repositories.NewConstraintRepo(connpool)
	constraintService := services.ConstraintService{Repository: constraintRepository}
	constraintController := ConstraintController{Service: constraintService}

	fmt.Println("Constraint controller up and running.")

	return constraintController
}

func GetExternalRightController(connpool *pgxpool.Pool) ExternalRightController {
	externalRightRepository := repositories.NewExternalRightRepo(connpool)
	externalRightService := services.ExternalRightService{Repository: externalRightRepository}
	externalRightController := ExternalRightController{Service: externalRightService}

	fmt.Println("ExternalRight controller up and running.")

	return externalRightController
}

func GetEmployeeController(connpool *pgxpool.Pool) EmployeeController {
	employeeRepository := repositories.NewEmployeeRepo(connpool)
	employeeService := services.EmployeeService{Repository: employeeRepository}
	employeeController := EmployeeController{Service: employeeService}

	fmt.Println("Employee controller up and running.")

	return employeeController
}

func getProductController(connpool *pgxpool.Pool, employeeRepo *repositories.EmployeeRepository) ProductController {

	productRepository := repositories.NewProductRepo(connpool)
	productService := services.ProductService{ProductRepository: productRepository, EmployeeRepository: *employeeRepo}
	productController := ProductController{Service: productService}

	fmt.Println("Product controller up and running.")

	return productController
}


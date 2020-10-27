package controllers

import (
	"context"
	"fmt"
	"internship_project/models"
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

	testAdmin     models.Employee
	testEmployee  models.Employee
	testEmployee1 models.Employee

	testProduct  models.Product
	testProduct1 models.Product

	testCompany  models.Company
	testCompany1 models.Company
	testCompany2 models.Company
	mainCompany1 models.Company

	testEar  models.ExternalRights
	testEar1 models.ExternalRights
	testEar2 models.ExternalRights
	testEar3 models.ExternalRights

	testConstraint  models.AccessConstraint
	testConstraint2 models.AccessConstraint
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
	ProductCont = GetProductController(connpool)
	EmployeeCont = GetEmployeeController(connpool)
	ConstraintCont = GetConstraintController(connpool)
	ExternalRightCont = GetExternalRightController(connpool)

	testCompany = models.Company{
		Id:     "",
		Name:   "SpaceX",
		IsMain: false,
	}

	testCompany1 = models.Company{
		Id:     "153fac6d-760d-4841-87e9-15aee2f25182",
		Name:   "Test Kompanija 1",
		IsMain: false,
	}

	testCompany2 = models.Company{
		Id:     "aac4ca5c-315f-4c16-8e36-eb62e0292d25",
		Name:   "Test Kompanija 2",
		IsMain: false,
	}

	mainCompany1 = models.Company{
		Id:     "91f88893-5cef-4d3c-9d6a-ed120f7e449e",
		Name:   "Main Kompanija 1",
		IsMain: true,
	}

	testEmployee = models.Employee{
		ID:        "",
		FirstName: "Test Name",
		LastName:  "Test Surname",
		CompanyID: testCompany1.Id,
		C:         false,
		R:         true,
		U:         true,
		D:         true,
	}

	testEmployee1 = models.Employee{
		ID:        "3f17c2bb-d65c-4ac5-aadd-1f3d933ae860",
		FirstName: "Preadded",
		LastName:  "Test Surname",
		CompanyID: testCompany2.Id,
		C:         false,
		R:         true,
		U:         true,
		D:         true,
	}

	testAdmin = models.Employee{
		ID:        "9d6ffd16-89e1-4ece-9e7c-09d4bf390838",
		FirstName: "Admin",
		LastName:  "Admin",
		CompanyID: testCompany1.Id,
		C:         true,
		R:         true,
		U:         true,
		D:         true,
	}

	testProduct = models.Product{
		ID:       "",
		Name:     "TEST_PRODUCT",
		Price:    99,
		Quantity: 10,
		IDC:      testCompany1.Id,
	}

	testProduct1 = models.Product{
		ID:       "a8451090-9e22-4fc2-832b-c65d0fc080c8",
		Name:     "Inserted Product",
		Price:    99,
		Quantity: 10,
		IDC:      testCompany1.Id,
	}

	testEar = models.ExternalRights{
		ID:       "a3a3d913-12d6-444b-aa21-ed1eb33bbde2",
		Read:     true,
		Update:   false,
		Delete:   false,
		Approved: false,
		IDSC:     testCompany1.Id,
		IDRC:     testCompany2.Id,
	}

	testEar1 = models.ExternalRights{
		ID:       "6b64bc14-01c5-4afa-8ff9-40545b8d0939",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: true,
		IDSC:     testCompany1.Id,
		IDRC:     testCompany2.Id,
	}

	testEar2 = models.ExternalRights{
		ID:       "c9f8384b-1615-4117-a983-00d574c2614c",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: false,
		IDSC:     testCompany2.Id,
		IDRC:     testCompany1.Id,
	}

	testEar3 = models.ExternalRights{
		ID:       "0fa90c22-18a4-4ed8-a834-190b1b84e32e",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: false,
		IDSC:     testCompany1.Id,
		IDRC:     testCompany2.Id,
	}

	testConstraint = models.AccessConstraint{
		ID:            "0657e510-1775-45be-a6cc-6fff7a9841fc",
		IDEAR:         testEar1.ID,
		OperatorID:    2,
		PropertyID:    1,
		PropertyValue: 15,
	}

	testConstraint2 = models.AccessConstraint{
		ID:            "d70583a4-d5a3-4d6e-847f-dfba98bd3a27",
		IDEAR:         testEar1.ID,
		OperatorID:    3,
		PropertyID:    1,
		PropertyValue: 20,
	}

	SetUpTables(connpool)

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
	companyRepository := repositories.CompanyRepository{DB: connpool}
	companyService := services.CompanyService{Repository: companyRepository}
	companyController := CompanyController{Service: companyService}

	fmt.Println("Company controller up and running.")

	return companyController
}

func GetConstraintController(connpool *pgxpool.Pool) ConstraintController {
	constraintRepository := repositories.ConstraintRepository{DB: connpool}
	constraintService := services.ConstraintService{Repository: constraintRepository}
	constraintController := ConstraintController{Service: constraintService}

	fmt.Println("Constraint controller up and running.")

	return constraintController
}

func GetExternalRightController(connpool *pgxpool.Pool) ExternalRightController {
	externalRightRepository := repositories.ExternalRightRepository{DB: connpool}
	externalRightService := services.ExternalRightService{Repository: externalRightRepository}
	externalRightController := ExternalRightController{Service: externalRightService}

	fmt.Println("ExternalRight controller up and running.")

	return externalRightController
}

func GetEmployeeController(connpool *pgxpool.Pool) EmployeeController {
	employeeRepository := repositories.EmployeeRepository{DB: connpool}
	employeeService := services.EmployeeService{Repository: employeeRepository}
	employeeController := EmployeeController{Service: employeeService}

	fmt.Println("Employee controller up and running.")

	return employeeController
}

func GetProductController(connpool *pgxpool.Pool) ProductController {
	productRepository := repositories.ProductRepository{DB: connpool}
	employeeRepository := repositories.EmployeeRepository{DB: connpool}
	productService := services.ProductService{ProductRepository: productRepository, EmployeeRepository: employeeRepository}
	productController := ProductController{Service: productService}

	fmt.Println("Product controller up and running.")

	return productController
}

func insertMockData(db *pgxpool.Pool) {
	// Insert Companies
	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany1.Id, testCompany1.Name, testCompany1.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany2.Id, testCompany2.Name, testCompany2.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		mainCompany1.Id, mainCompany1.Name, mainCompany1.IsMain)

	// Insert Product
	db.Exec(context.Background(), "insert into products (id, name, price, quantity, idc) VALUES($1, $2, $3, $4, $5)",
		testProduct1.ID, testProduct1.Name, testProduct1.Price, testProduct1.Quantity, testProduct1.IDC)

	// Insert Users
	db.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		testAdmin.ID, testAdmin.FirstName, testAdmin.LastName, testAdmin.CompanyID, testAdmin.C, testAdmin.R, testAdmin.U, testAdmin.D)

	db.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		testEmployee1.ID, testEmployee1.FirstName, testEmployee1.LastName, testEmployee1.CompanyID, testEmployee1.C, testEmployee1.R, testEmployee1.U, testEmployee1.D)

	// Insert external access rights
	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		testEar1.ID, testEar1.IDSC, testEar1.IDRC, testEar1.Read, testEar1.Update, testEar1.Delete, testEar1.Approved)

	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		testEar2.ID, testEar2.IDSC, testEar2.IDRC, testEar2.Read, testEar2.Update, testEar2.Delete, testEar2.Approved)

	// Insert Properties
	db.Exec(context.Background(), "insert into properties (id, name) values ($1, $2)",
		"1", "quantity")

	db.Exec(context.Background(), "insert into properies (id, name) values ($1, $2)",
		"2", "price")

	// Insert Operators
	db.Exec(context.Background(), "insert into operators (id, name) values ($1, $2)",
		"1", ">")

	db.Exec(context.Background(), "insert into operators (id, name) values ($1, $2)",
		"2", ">=")

	db.Exec(context.Background(), "insert into operators (id, name) values ($1, $2)",
		"3", "<")

	db.Exec(context.Background(), "insert into operators (id, name) values ($1, $2)",
		"4", "<=")

	// Insert Constraints
	db.Exec(context.Background(), "insert into access_constraints (id, idear, operator_id, property_id, property_value) VALUES($1, $2, $3, $4, $5)",
		testConstraint2.ID, testConstraint2.IDEAR, testConstraint2.OperatorID, testConstraint2.PropertyID, testConstraint2.PropertyValue)
}

func SetUpTables(db *pgxpool.Pool) {
	utils.DropTables(db)
	utils.CreateTables(db)
	insertMockData(db)
}

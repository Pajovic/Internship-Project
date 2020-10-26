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
	connpool    *pgxpool.Pool
	CompanyCont CompanyController

	testCompany models.Company

	testCompany1 models.Company
	testCompany2 models.Company
	mainCompany1 models.Company
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

func insertMockData(db *pgxpool.Pool) {
	// Insert Companies
	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany1.Id, testCompany1.Name, testCompany1.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany2.Id, testCompany2.Name, testCompany2.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		mainCompany1.Id, mainCompany1.Name, mainCompany1.IsMain)
}

func SetUpTables(db *pgxpool.Pool) {
	utils.DropTables(db)
	utils.CreateTables(db)
	insertMockData(db)
}

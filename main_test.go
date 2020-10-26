package main

import (
	"context"
	"fmt"
	"internship_project/controllers"
	"internship_project/models"
	"internship_project/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
	"github.com/stretchr/testify/assert"
)

var (
	connpool          *pgxpool.Pool
	CompanyController controllers.CompanyController

	testCompany models.Company

	testCompany1 models.Company
	testCompany2 models.Company
	mainCompany1 models.Company
)

func TestMain(m *testing.M) {
	connpool = GetTestConnectionPool()
	defer connpool.Close()

	CompanyController = GetCompanyController(connpool)

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

	os.Exit(m.Run())
}

func GetTestConnectionPool() *pgxpool.Pool {
	var conf Config
	if _, err := confl.DecodeFile("dbconfig.conf", &conf); err != nil {
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

}

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/company", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(CompanyController.GetAllCompanies)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(rr.Code, 400, "Response code is not 500")
		assert.Equal(rr.Body.String(), `relation "public.companies" does not exist`, "Error message is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(rr.Code, http.StatusOK, "Response code is not 200")
	})

}

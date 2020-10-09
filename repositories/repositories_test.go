package repositories

import (
	"context"
	"fmt"
	"internship_project/models"
	"os"
	"testing"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
	uuid "github.com/satori/go.uuid"
)

type config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

var EmployeeRepo EmployeeRepository
var ProductRepo ProductRepository
var CompanyRepo CompanyRepository

var testEmployee models.Employee
var testProduct models.Product
var testCompany models.Company

func TestMain(m *testing.M) {
	connpool := getConnPool()
	defer connpool.Close()

	EmployeeRepo = EmployeeRepository{DB: connpool}
	ProductRepo = ProductRepository{DB: connpool}
	CompanyRepo = CompanyRepository{DB: connpool}

	testCompany = models.Company{
		Id:     "",
		Name:   "SpaceX",
		IsMain: false,
	}

	testEmployee = models.Employee{
		ID:        "",
		FirstName: "Test Name",
		LastName:  "Test Surname",
		CompanyID: "153fac6d-760d-4841-87e9-15aee2f25182", // ID from database
		C:         false,
		R:         true,
		U:         false,
		D:         false,
	}

	testProduct = models.Product{
		ID:       "",
		Name:     "TEST_PRODUCT",
		Price:    99,
		Quantity: 10,
		IDC:      "153fac6d-760d-4841-87e9-15aee2f25182",
	}

	testCleanup(connpool)
	defer testCleanup(connpool)

	setupTables(connpool)

	code := m.Run()

	os.Exit(code)
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

func DoesTableExist(tableName string, connpool *pgxpool.Pool) bool {
	var n int64
	err := connpool.QueryRow(context.Background(), "select 1 from information_schema.tables where table_name=$1", tableName).Scan(&n)
	if err == pgx.ErrNoRows || err != nil {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

func getConnPool() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("./../dbconfig.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.TestDatabaseURL)

	dbtest, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbtest
}

func setupTables(db *pgxpool.Pool) {
	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		"153fac6d-760d-4841-87e9-15aee2f25182", "Test Kompanija", true)
}

func testCleanup(db *pgxpool.Pool) {
	db.Exec(context.Background(), "DELETE FROM products;")
	db.Exec(context.Background(), "DELETE FROM employees;")
	db.Exec(context.Background(), "DELETE FROM companies;")
}

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
var EarRepo EarRepository
var ConstraintRepo ConstraintRepository

var testAdmin models.Employee

var testEmployee models.Employee
var testProduct models.Product
var testCompany models.Company
var testEar models.ExternalRights

var testCompany1 models.Company
var testCompany2 models.Company
var mainCompany1 models.Company

var testEar1 models.ExternalRights
var testEar2 models.ExternalRights

var testConstraint models.AccessConstraint

func TestMain(m *testing.M) {
	connpool := getConnPool()
	defer connpool.Close()

	EmployeeRepo = EmployeeRepository{DB: connpool}
	ProductRepo = ProductRepository{DB: connpool}
	CompanyRepo = CompanyRepository{DB: connpool}
	EarRepo = EarRepository{DB: connpool}
	ConstraintRepo = ConstraintRepository{DB: connpool}

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
		IDSC:     "",
		IDRC:     "",
	}

	testEar2 = models.ExternalRights{
		ID:       "c9f8384b-1615-4117-a983-00d574c2614c",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: false,
		IDSC:     "",
		IDRC:     "",
	}

	testConstraint = models.AccessConstraint{
		ID:            "",
		IDEAR:         testEar1.ID,
		OperatorID:    2,
		PropertyID:    1,
		PropertyValue: 15,
	}

	SetupTables(connpool)

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
	if _, configerr := confl.DecodeFile("./../dbconfig.conf", &conf); configerr != nil {
		panic(configerr)
	}

	poolConfig, poolerr := pgxpool.ParseConfig(conf.TestDatabaseURL)
	if poolerr != nil {
		panic("Error configuring pool")
	}

	dbtest, dberr := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if dberr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dberr)
		os.Exit(1)
	}

	return dbtest
}

func insertMockData(db *pgxpool.Pool) {
	// Insert Companies
	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany1.Id, testCompany1.Name, testCompany1.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		testCompany2.Id, testCompany2.Name, testCompany2.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		mainCompany1.Id, mainCompany1.Name, mainCompany1.IsMain)

	// Insert Admin
	db.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		testAdmin.ID, testAdmin.FirstName, testAdmin.LastName, testAdmin.CompanyID, testAdmin.C, testAdmin.R, testAdmin.U, testAdmin.D)

	// Insert external access rights
	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		testEar1.ID, testCompany1.Id, testCompany2.Id, testEar1.Read, testEar1.Update, testEar1.Delete, testEar1.Approved)

	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		testEar2.ID, testCompany2.Id, testCompany1.Id, testEar2.Read, testEar2.Update, testEar2.Delete, testEar2.Approved)

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
}

func SetupTables(db *pgxpool.Pool) {
	DropTables(db)
	CreateTables(db)
}

func CreateTables(db *pgxpool.Pool) {
	// Companies
	db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS companies (
		id uuid NOT NULL,
		"name" varchar(30) NOT NULL,
		ismain bool NOT NULL,
		CONSTRAINT companies_pk PRIMARY KEY (id)
	);`)

	// Employees
	db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS employees (
		id uuid NOT NULL,
		firstname varchar(30) NOT NULL,
		lastname varchar(30) NOT NULL,
		idc uuid NOT NULL,
		c bool NOT NULL,
		r bool NOT NULL,
		u bool NOT NULL,
		d bool NOT NULL,
		CONSTRAINT employees_pk PRIMARY KEY (id)
	);
	ALTER TABLE employees ADD CONSTRAINT employees_fk FOREIGN KEY (idc) REFERENCES companies(id);
	`)

	// Products
	db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS products (
		id uuid NOT NULL,
		"name" varchar(30) NOT NULL,
		price float4 NOT NULL,
		quantity int4 NOT NULL,
		idc uuid NOT NULL,
		CONSTRAINT products_pk PRIMARY KEY (id)
	);
	ALTER TABLE products ADD CONSTRAINT products_fk FOREIGN KEY (idc) REFERENCES companies(id);
	`)

	// External Access Rights
	db.Exec(context.Background(), ` 
	CREATE TABLE IF NOT EXISTS external_access_rights (
		id uuid NOT NULL,
		idsc uuid NOT NULL,
		idrc uuid NOT NULL,
		r bool NOT NULL,
		u bool NOT NULL,
		d bool NOT NULL,
		approved bool NOT NULL,
		CONSTRAINT external_access_rights_pk PRIMARY KEY (id)
	);

	ALTER TABLE external_access_rights ADD CONSTRAINT external_access_rights_idrc FOREIGN KEY (idrc) REFERENCES companies(id);
	ALTER TABLE external_access_rights ADD CONSTRAINT external_access_rights_idsc FOREIGN KEY (idsc) REFERENCES companies(id);
	`)

	// Access Constraints
	db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS operators (
			id int4 NOT NULL,
			"name" varchar(5) NOT NULL,
			CONSTRAINT operators_pk PRIMARY KEY (id)
		);

		CREATE TABLE IF NOT EXISTS properties (
			id int8 NOT NULL,
			"name" varchar(20) NOT NULL,
			CONSTRAINT properties_pk PRIMARY KEY (id)
		);

		CREATE TABLE IF NOT EXISTS access_constraints (
			id uuid NOT NULL,
			idear uuid NOT NULL,
			operator_id int4 NOT NULL,
			property_id int8 NOT NULL,
			property_value float8 NOT NULL,
			CONSTRAINT access_constraints_pk PRIMARY KEY (id)
		);

		ALTER TABLE access_constraints ADD CONSTRAINT access_constraints_idear FOREIGN KEY (idear) REFERENCES external_access_rights(id);
		ALTER TABLE access_constraints ADD CONSTRAINT access_constraints_operator_id FOREIGN KEY (operator_id) REFERENCES operators(id);
		ALTER TABLE access_constraints ADD CONSTRAINT access_constraints_property_id FOREIGN KEY (property_id) REFERENCES properties(id);
	`)

	insertMockData(db)
}

func DropTables(db *pgxpool.Pool) {
	db.Exec(context.Background(), "DROP TABLE IF EXISTS products;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS employees;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS access_constraints;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS operators;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS properties;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS external_access_rights;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS companies;")
}

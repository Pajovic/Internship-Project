package repositories

import (
	"context"
	"fmt"
	"internship_project/models"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
)

type config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

var DbInstance *pgxpool.Pool

var testEmployee models.Employee

func TestMain(m *testing.M) {
	DbInstance = instantiateDb()
	testEmployee = models.Employee{
		ID:        "517a9da1-ba1a-492f-8df2-e695df582bf9",
		FirstName: "Test Name",
		LastName:  "Test Surname",
		CompanyID: "153fac6d-760d-4841-87e9-15aee2f25182", // ID from database
		C:         false,
		R:         true,
		U:         false,
		D:         false,
	}

	testCleanup(DbInstance)
	defer testCleanup(DbInstance)

	setupTables(DbInstance)

	code := m.Run()

	os.Exit(code)
}

func instantiateDb() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("./../database.conf", &conf); err != nil {
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
	db.Exec(context.Background(), `CREATE TABLE public.companies (
		id uuid NOT NULL,
		"name" varchar(30) NOT NULL,
		ismain bool NOT NULL,
		CONSTRAINT companies_pk PRIMARY KEY (id)
	);`)

	db.Exec(context.Background(), `CREATE TABLE public.employees (
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
	ALTER TABLE public.employees ADD CONSTRAINT employees_fk FOREIGN KEY (idc) REFERENCES companies(id);
	`)

	db.Exec(context.Background(), `CREATE TABLE public.products (
		id uuid NOT NULL,
		"name" varchar(30) NOT NULL,
		price float4 NOT NULL,
		quantity int4 NOT NULL,
		idc uuid NOT NULL,
		CONSTRAINT products_pk PRIMARY KEY (id)
	);
	ALTER TABLE public.products ADD CONSTRAINT products_fk FOREIGN KEY (idc) REFERENCES companies(id);`)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		"153fac6d-760d-4841-87e9-15aee2f25182", "Test Kompanija", true)
}

func testCleanup(db *pgxpool.Pool) {
	db.Exec(context.Background(), "DROP TABLE IF EXISTS products;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS employees;")
	db.Exec(context.Background(), "DROP TABLE IF EXISTS companies;")
}

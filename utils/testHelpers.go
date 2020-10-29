package utils

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"internship_project/models"
)

var (
	TestAdmin models.Employee = models.Employee{
		ID:        "9d6ffd16-89e1-4ece-9e7c-09d4bf390838",
		FirstName: "Admin",
		LastName:  "Admin",
		CompanyID: TestCompany1.ID,
		C:         true,
		R:         true,
		U:         true,
		D:         true,
	}
	TestEmployee models.Employee = models.Employee{
		ID:        "298dd516-9663-4b96-bc3c-8e0b2b9be469",
		FirstName: "Test Name",
		LastName:  "Test Surname",
		CompanyID: TestCompany1.ID,
		C:         false,
		R:         true,
		U:         true,
		D:         true,
	}
	TestEmployee1 models.Employee = models.Employee{
		ID:        "3f17c2bb-d65c-4ac5-aadd-1f3d933ae860",
		FirstName: "Preadded",
		LastName:  "Test Surname",
		CompanyID: TestCompany2.ID,
		C:         false,
		R:         true,
		U:         true,
		D:         true,
	}

	TestProduct  models.Product = models.Product{
		ID:       "",
		Name:     "TEST_PRODUCT",
		Price:    99,
		Quantity: 10,
		IDC:      TestCompany1.ID,
	}
	TestProduct1 models.Product = models.Product{
		ID:       "a8451090-9e22-4fc2-832b-c65d0fc080c8",
		Name:     "Inserted Product",
		Price:    99,
		Quantity: 10,
		IDC:      TestCompany1.ID,
	}

	TestCompany  models.Company = models.Company{
		ID:     "",
		Name:   "SpaceX",
		IsMain: false,
	}
	TestCompany1 models.Company = models.Company{
		ID:     "153fac6d-760d-4841-87e9-15aee2f25182",
		Name:   "Test Kompanija 1",
		IsMain: false,
	}
	TestCompany2 models.Company = models.Company{
		ID:     "aac4ca5c-315f-4c16-8e36-eb62e0292d25",
		Name:   "Test Kompanija 2",
		IsMain: false,
	}
	MainCompany1 models.Company = models.Company{
		ID:     "91f88893-5cef-4d3c-9d6a-ed120f7e449e",
		Name:   "Main Kompanija 1",
		IsMain: true,
	}

	TestEar models.ExternalRights = models.ExternalRights{
		ID:       "a3a3d913-12d6-444b-aa21-ed1eb33bbde2",
		Read:     true,
		Update:   false,
		Delete:   false,
		Approved: false,
		IDSC:     TestCompany1.ID,
		IDRC:     TestCompany2.ID,
	}
	TestEar1 models.ExternalRights = models.ExternalRights{
		ID:       "6b64bc14-01c5-4afa-8ff9-40545b8d0939",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: true,
		IDSC:     TestCompany1.ID,
		IDRC:     TestCompany2.ID,
	}
	TestEar2 models.ExternalRights = models.ExternalRights{
		ID:       "c9f8384b-1615-4117-a983-00d574c2614c",
		Read:     true,
		Update:   true,
		Delete:   true,
		Approved: false,
		IDSC:     TestCompany2.ID,
		IDRC:     TestCompany1.ID,
	}

	TestConstraint models.AccessConstraint = models.AccessConstraint{
		ID:            "",
		IDEAR:         TestEar1.ID,
		OperatorID:    2,
		PropertyID:    1,
		PropertyValue: 15,
	}
	TestConstraint2 models.AccessConstraint = models.AccessConstraint{
		ID:            "d70583a4-d5a3-4d6e-847f-dfba98bd3a27",
		IDEAR:         TestEar1.ID,
		OperatorID:    3,
		PropertyID:    1,
		PropertyValue: 20,
	}
)

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

func insertMockData(db *pgxpool.Pool) {
	// Insert Companies
	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		TestCompany1.ID, TestCompany1.Name, TestCompany1.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		TestCompany2.ID, TestCompany2.Name, TestCompany2.IsMain)

	db.Exec(context.Background(), "insert into companies (id, name, ismain) values ($1, $2, $3)",
		MainCompany1.ID, MainCompany1.Name, MainCompany1.IsMain)

	// Insert Product
	db.Exec(context.Background(), "insert into products (id, name, price, quantity, idc) VALUES($1, $2, $3, $4, $5)",
		TestProduct1.ID, TestProduct1.Name, TestProduct1.Price, TestProduct1.Quantity, TestProduct1.IDC)

	// Insert Users
	db.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		TestAdmin.ID, TestAdmin.FirstName, TestAdmin.LastName, TestAdmin.CompanyID, TestAdmin.C, TestAdmin.R, TestAdmin.U, TestAdmin.D)

	db.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		TestEmployee1.ID, TestEmployee1.FirstName, TestEmployee1.LastName, TestEmployee1.CompanyID, TestEmployee1.C, TestEmployee1.R, TestEmployee1.U, TestEmployee1.D)

	// Insert external access rights
	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		TestEar1.ID, TestEar1.IDSC, TestEar1.IDRC, TestEar1.Read, TestEar1.Update, TestEar1.Delete, TestEar1.Approved)

	db.Exec(context.Background(), `insert into external_access_rights (id, idsc, idrc, r, u, d, approved) values ($1, $2, $3, $4, $5, $6, $7)`,
		TestEar2.ID, TestEar2.IDSC, TestEar2.IDRC, TestEar2.Read, TestEar2.Update, TestEar2.Delete, TestEar2.Approved)

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
		TestConstraint2.ID, TestConstraint2.IDEAR, TestConstraint2.OperatorID, TestConstraint2.PropertyID, TestConstraint2.PropertyValue)
}

func SetUpTables(db *pgxpool.Pool) {
	DropTables(db)
	CreateTables(db)
	insertMockData(db)
}

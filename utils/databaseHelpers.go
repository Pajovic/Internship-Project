package utils

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
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

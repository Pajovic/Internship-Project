-- public.companies definition

CREATE TABLE public.companies (
	id uuid NOT NULL,
	"name" varchar(30) NOT NULL,
	ismain bool NOT NULL,
	CONSTRAINT companies_pk PRIMARY KEY (id)
);

-- public.employees definition

CREATE TABLE public.employees (
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


-- public.employees foreign keys

ALTER TABLE public.employees ADD CONSTRAINT employees_fk FOREIGN KEY (idc) REFERENCES companies(id);

-- public.products definition

CREATE TABLE public.products (
	id uuid NOT NULL,
	"name" varchar(30) NOT NULL,
	price float4 NOT NULL,
	quantity int4 NOT NULL,
	idc uuid NOT NULL,
	CONSTRAINT products_pk PRIMARY KEY (id)
);

-- public.products foreign keys

ALTER TABLE public.products ADD CONSTRAINT products_fk FOREIGN KEY (idc) REFERENCES companies(id);

-- public.operators definition

CREATE TABLE public.operators (
	id int4 NOT NULL,
	"name" varchar(5) NOT NULL,
	CONSTRAINT operators_pk PRIMARY KEY (id)
);

-- public.properties definition

CREATE TABLE public.properties (
	id int8 NOT NULL,
	"name" varchar(20) NOT NULL,
	CONSTRAINT properties_pk PRIMARY KEY (id)
);

-- public.external_access_rights definition

CREATE TABLE public.external_access_rights (
	id uuid NOT NULL,
	idsc uuid NOT NULL,
	idrc uuid NOT NULL,
	r bool NOT NULL,
	u bool NOT NULL,
	d bool NOT NULL,
	approved bool NOT NULL,
	CONSTRAINT external_access_rights_pk PRIMARY KEY (id)
);


-- public.external_access_rights foreign keys

ALTER TABLE public.external_access_rights ADD CONSTRAINT external_access_rights_idrc FOREIGN KEY (idrc) REFERENCES companies(id);
ALTER TABLE public.external_access_rights ADD CONSTRAINT external_access_rights_idsc FOREIGN KEY (idsc) REFERENCES companies(id);

-- public.access_constraints definition

CREATE TABLE public.access_constraints (
	id uuid NOT NULL,
	idear uuid NOT NULL,
	operator_id int4 NOT NULL,
	property_id int8 NOT NULL,
	property_value float8 NOT NULL,
	CONSTRAINT access_constraints_pk PRIMARY KEY (id)
);

-- public.access_constraints foreign keys

ALTER TABLE public.access_constraints ADD CONSTRAINT access_constraints_idear FOREIGN KEY (idear) REFERENCES external_access_rights(id);
ALTER TABLE public.access_constraints ADD CONSTRAINT access_constraints_operator_id FOREIGN KEY (operator_id) REFERENCES operators(id);
ALTER TABLE public.access_constraints ADD CONSTRAINT access_constraints_property_id FOREIGN KEY (property_id) REFERENCES properties(id);

-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id varchar NOT NULL,
	email varchar NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);


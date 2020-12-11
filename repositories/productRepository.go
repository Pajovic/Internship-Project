package repositories

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	GetAllProducts(string) ([]models.Product, error)
	GetProduct(string, string) (models.Product, error)
	AddProduct(*models.Product) error
	UpdateProduct(models.Product) error
	DeleteProduct(string) error
	DeleteProductsFromCompany(string) error
}

type productRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) ProductRepository {
	if db == nil {
		panic("ProductRepository not created, pgxpool is nil")
	}
	return &productRepository{
		DB: db,
	}
}

func (repository *productRepository) GetAllProducts(employeeIdc string) ([]models.Product, error) {
	earConstraints := []models.EarConstraint{}

	query := `select ear.id "idear", ear.idrc, ear.idsc, coalesce(p.name::varchar(20), '') as "property",
	coalesce(o2.name::varchar(5), '') as "operator", coalesce(ac.property_value::int4, 0)
    from external_access_rights ear left outer join access_constraints ac on ear.id = ac.idear
	left outer join operators o2 on o2.id = ac.operator_id 
	left outer join properties p on p.id = ac.property_id 
	where ear.idrc = $1 and ear.r = true and ear.approved = true;`

	rows, err := repository.DB.Query(context.Background(), query, employeeIdc)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var earConstraint models.EarConstraint
		err := rows.Scan(&earConstraint.IDEAR, &earConstraint.IDRC, &earConstraint.IDSC, &earConstraint.Property, &earConstraint.Operator, &earConstraint.PropertyValue)
		if err != nil {
			return nil, err
		}
		earConstraints = append(earConstraints, earConstraint)
	}

	finalQueryTemplate := `
	select * from products p where p.idc = $1
	{{- range . -}}	
		{{- if .Operator}}
			or (p.idc = '{{.IDSC}}' and p.{{.Property}} {{.Operator}} {{.PropertyValue}})
		{{- else}}
			or (p.idc = '{{.IDSC}}')
		{{- end -}}
	{{- end -}}
	;`

	var buff bytes.Buffer
	t := template.Must(template.New("getProducts").Parse(finalQueryTemplate))

	err = t.Execute(&buff, earConstraints)
	if err != nil {
		return nil, err
	}

	finalQuery := strings.TrimSpace(buff.String())
	fmt.Println(finalQuery)

	rowsProducts, err := repository.DB.Query(context.Background(), finalQuery, employeeIdc)
	defer rowsProducts.Close()

	if err != nil {
		return nil, err
	}

	products := []models.Product{}

	for rowsProducts.Next() {
		var productPers persistence.Products
		productPers.Scan(&rowsProducts)

		var productUUID string
		err = productPers.Id.AssignTo(&productUUID)
		if err != nil {
			return nil, err
		}

		var companyUUID string
		err = productPers.Idc.AssignTo(&companyUUID)
		if err != nil {
			return nil, err
		}

		product := models.Product{
			ID:       productUUID,
			Name:     productPers.Name,
			Price:    productPers.Price,
			Quantity: productPers.Quantity,
			IDC:      companyUUID,
		}

		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) GetProduct(id string, employeeIdc string) (models.Product, error) {
	product := models.Product{}
	earConstraints := []models.EarConstraint{}

	query := `select ear.id "idear", ear.idrc, ear.idsc, coalesce(p.name::varchar(20), '') as "property",
	coalesce(o2.name::varchar(5), '') as "operator", coalesce(ac.property_value::int4, 0)
    from external_access_rights ear left outer join access_constraints ac on ear.id = ac.idear
	left outer join operators o2 on o2.id = ac.operator_id 
	left outer join properties p on p.id = ac.property_id 
	where ear.idrc = $1 and ear.r = true and ear.approved = true;`

	rows, err := repository.DB.Query(context.Background(), query, employeeIdc)
	defer rows.Close()
	if err != nil {
		return product, err
	}

	for rows.Next() {
		var earConstraint models.EarConstraint
		err := rows.Scan(&earConstraint.IDEAR, &earConstraint.IDRC, &earConstraint.IDSC, &earConstraint.Property, &earConstraint.Operator, &earConstraint.PropertyValue)
		if err != nil {
			return product, err
		}
		earConstraints = append(earConstraints, earConstraint)
	}

	finalQueryTemplate := `
	select * from products p where p.idc = $1
	{{- range . -}}	
		{{- if .Operator}}
			or (p.idc = '{{.IDSC}}' and p.{{.Property}} {{.Operator}} {{.PropertyValue}})
		{{- else}}
			or (p.idc = '{{.IDSC}}')
		{{- end -}}
	{{- end -}}
	;`

	var buff bytes.Buffer
	t := template.Must(template.New("getProducts").Parse(finalQueryTemplate))

	err = t.Execute(&buff, earConstraints)
	if err != nil {
		return product, err
	}

	finalQuery := strings.TrimSpace(buff.String())
	fmt.Println(finalQuery)

	rowsProducts, err := repository.DB.Query(context.Background(), finalQuery, employeeIdc)
	defer rowsProducts.Close()

	if err != nil {
		return product, err
	}

	if !rowsProducts.Next() {
		return product, errors.New("There is no product with this ID")
	}

	for rowsProducts.Next() {
		var productPers persistence.Products
		productPers.Scan(&rowsProducts)

		var productUUID string
		err = productPers.Id.AssignTo(&productUUID)
		if err != nil {
			return product, err
		}

		var companyUUID string
		err = productPers.Idc.AssignTo(&companyUUID)
		if err != nil {
			return product, err
		}

		product = models.Product{
			ID:       productUUID,
			Name:     productPers.Name,
			Price:    productPers.Price,
			Quantity: productPers.Quantity,
			IDC:      companyUUID,
		}
	}

	return product, nil
}

func (repository *productRepository) AddProduct(product *models.Product) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	product.ID = uuid.NewV4().String()

	productPers := persistence.Products{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
	productPers.Idc.Set(product.IDC)
	productPers.Id.Set(product.ID)

	_, err = productPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *productRepository) UpdateProduct(product models.Product) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	productPers := persistence.Products{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
	productPers.Idc.Set(product.IDC)
	productPers.Id.Set(product.ID)

	commandTag, err := productPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *productRepository) DeleteProduct(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	productPers := persistence.Products{}
	productPers.Id.Set(id)

	commandTag, err := productPers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *productRepository) DeleteProductsFromCompany(idc string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `DELETE FROM products WHERE idc=$1`

	_, err = tx.Exec(context.Background(), query, idc)

	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

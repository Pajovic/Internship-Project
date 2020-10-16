package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func (repository *ProductRepository) GetAllProducts(employeeIdc string) ([]models.Product, error) {
	earConstraints := []models.EarConstraint{}

	query := `select ear.id "idear", ear.idrc, ear.idsc, p.name "property", o2.name "operator", ac.property_value from external_access_rights ear
	left outer join access_constraints ac on ear.id = ac.idear
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
		err := rows.Scan(&earConstraint.Idear, &earConstraint.Idrc, &earConstraint.Idsc, &earConstraint.Property, &earConstraint.Operator, &earConstraint.PropertyValue)
		if err != nil {
			return nil, err
		}
		earConstraints = append(earConstraints, earConstraint)
	}

	var finalQuery strings.Builder

	finalQuery.WriteString("select * from products p where p.idc = $1")

	for _, earc := range earConstraints {
		finalQuery.WriteString(" union select * from products p where p.idc = '" + earc.Idsc + "' ")
		if earc.Operator != "" && earc.Property != "" {
			finalQuery.WriteString("and p." + earc.Property + earc.Operator + strconv.Itoa(earc.PropertyValue))
		}
	}
	finalQuery.WriteString(";")

	products := []models.Product{}

	rowsProducts, err := repository.DB.Query(context.Background(), finalQuery.String(), employeeIdc)
	defer rowsProducts.Close()
	if err != nil {
		return nil, err
	}
	for rowsProducts.Next() {
		var product models.Product
		err := rowsProducts.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.IDC)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *ProductRepository) GetProduct(id string) (models.Product, error) {
	var product models.Product
	rows, err := repository.DB.Query(context.Background(), "select * from products where id=$1", id)
	defer rows.Close()
	if err != nil {
		return product, err
	}

	for rows.Next() {
		var productPers persistence.Products
		productPers.Scan(&rows)

		var productUUID string
		err := productPers.Id.AssignTo(&productUUID)
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

		break
	}

	return product, nil
}

func (repository *ProductRepository) AddProduct(product *models.Product) error {
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
	productPers.Id.Set(uuid.NewV4())

	_, err = productPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil
}

func (repository *ProductRepository) UpdateProduct(product models.Product) error {
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
		return errors.New("No row found to update")
	}

	tx.Commit(context.Background())
	return nil
}

func (repository *ProductRepository) DeleteProduct(id string) error {
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
		return errors.New("No row found to delete")
	}
	return nil
}

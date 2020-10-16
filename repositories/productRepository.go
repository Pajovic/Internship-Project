package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"strconv"

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

	finalQuery := "select * from products p where p.idc = $1"

	for _, earc := range earConstraints {
		finalQuery += " union select * from products p where p.idc = '" + earc.Idsc + "' "
		if earc.Operator != "" && earc.Property != "" {
			finalQuery += "and p." + earc.Property + earc.Operator + strconv.Itoa(earc.PropertyValue)
		}
	}
	finalQuery += ";"

	products := []models.Product{}

	rowsProducts, err := repository.DB.Query(context.Background(), finalQuery, employeeIdc)
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
	err := repository.DB.QueryRow(context.Background(),
		"SELECT * FROM products WHERE id=$1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.IDC)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (repository *ProductRepository) AddProduct(product *models.Product) error {
	u := uuid.NewV4()
	product.ID = u.String()
	_, err := repository.DB.Exec(context.Background(),
		"INSERT INTO products (id, name, price, quantity, idc) values ($1, $2, $3, $4, $5)",
		u.Bytes(), product.Name, product.Price, product.Quantity, product.IDC)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ProductRepository) UpdateProduct(product models.Product) error {
	commandTag, err := repository.DB.Exec(context.Background(),
		"UPDATE products SET name=$1, price=$2, quantity=$3, idc=$4 WHERE id=$5",
		product.Name, product.Price, product.Quantity, product.IDC, product.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

func (repository *ProductRepository) DeleteProduct(id string) error {
	commandTag, err := repository.DB.Exec(context.Background(), "DELETE FROM products WHERE id=$1;", id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

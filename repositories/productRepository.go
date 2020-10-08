package repositories

import (
	"context"
	"errors"
	"internship_project/models"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func (repository *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product = []models.Product{}
	rows, err := repository.DB.Query(context.Background(), "SELECT * FROM products")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.IDC)
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

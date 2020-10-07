package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	assert := assert.New(t)

	oldProducts, _ := repository.GetAllProducts()
	err := repository.AddProduct(&testProduct)
	newProducts, _ := repository.GetAllProducts()

	assert.NoError(err)
	assert.Equal(len(newProducts)-len(oldProducts), 1)
}

func TestGetAllProducts(t *testing.T) {
	assert := assert.New(t)
	allProducts, err := repository.GetAllProducts()

	assert.NoError(err)
	assert.NotNil(allProducts)
	assert.IsType(allProducts, []models.Product{})
}

func TestGetProduct(t *testing.T) {
	assert := assert.New(t)
	testID := testProduct.ID

	product, err := repository.GetProduct(testID)

	assert.NoError(err)
	assert.NotNil(product)
	assert.NotEmpty(product)

	assert.Equal(testID, product.ID)
}

func TestUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	productForUpdate, _ := repository.GetProduct(testProduct.ID)
	productForUpdate.Name = "UPDATED Name"

	err := repository.UpdateProduct(productForUpdate)

	assert.NoError(err)
}

func TestDeleteProduct(t *testing.T) {
	assert := assert.New(t)

	err := repository.DeleteProduct(testProduct.ID)

	assert.NoError(err)
}

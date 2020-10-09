package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	assert := assert.New(t)

	assert.True(DoesTableExist("products", ProductRepo.DB), "Table does not exist.")

	oldProducts, _ := ProductRepo.GetAllProducts()
	err := ProductRepo.AddProduct(&testProduct)
	newProducts, _ := ProductRepo.GetAllProducts()

	assert.NoError(err)
	assert.Equal(len(newProducts)-len(oldProducts), 1, "Product was not added.")
}

func TestGetAllProducts(t *testing.T) {
	assert := assert.New(t)

	allProducts, err := ProductRepo.GetAllProducts()

	assert.NoError(err)
	assert.NotNil(allProducts, "Products were nil.")
	assert.IsType(allProducts, []models.Product{})
}

func TestGetProduct(t *testing.T) {
	assert := assert.New(t)
	testID := testProduct.ID

	assert.True(IsValidUUID(testID), "Product ID is not valid.")

	product, err := ProductRepo.GetProduct(testID)

	assert.NoError(err)
	assert.NotNil(product, "Product is nil")
	assert.NotEmpty(product, "ID does not exist and the object is empty.") // ID does not exist

	assert.Equal(testID, product.ID, "Product ID and test ID do not match.")
}

func TestUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testProduct.ID), "Product ID is not valid.")

	productForUpdate, _ := ProductRepo.GetProduct(testProduct.ID)
	productForUpdate.Name = "UPDATED Name"

	err := ProductRepo.UpdateProduct(productForUpdate)

	assert.NoError(err, "Product was not updated.")
}

func TestDeleteProduct(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testProduct.ID), "Product ID is not valid.")

	err := ProductRepo.DeleteProduct(testProduct.ID)

	assert.NoError(err, "Product ID was not deleted.")
}

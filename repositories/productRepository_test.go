package repositories

import (
	"internship_project/models"
	"internship_project/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		assert.False(DoesTableExist("products", Connpool))
		err := ProductRepo.AddProduct(&utils.TestProduct)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		oldProducts, _ := ProductRepo.GetAllProducts(utils.TestAdmin.CompanyID)
		err := ProductRepo.AddProduct(&utils.TestProduct)
		newProducts, _ := ProductRepo.GetAllProducts(utils.TestAdmin.CompanyID)
		assert.NoError(err)
		assert.Equal(1, len(newProducts)-len(oldProducts), "Product was not added.")
	})

	t.Run("add an existing product", func(t *testing.T) {
		existingProduct := &models.Product{ID: utils.TestProduct.ID}
		err := ProductRepo.AddProduct(existingProduct)

		assert.Error(err)
	})
}

func TestGetAllProducts(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful GetAll query", func(t *testing.T) {
		allProducts, err := ProductRepo.GetAllProducts(utils.TestAdmin.CompanyID)

		assert.NoError(err)
		assert.NotNil(allProducts, "Products were nil.")
		assert.IsType(allProducts, []models.Product{})
	})
}

func TestGetProduct(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		_, err := ProductRepo.GetProduct(invalidID)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "c5ef08c6-60eb-4687-bcbb-df37ebc9e105"
		assert.True(IsValidUUID(randomUUID))
		_, err := ProductRepo.GetProduct(randomUUID)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		testID := utils.TestProduct.ID
		product, err := ProductRepo.GetProduct(testID)

		assert.NoError(err)
		assert.NotNil(product, "Product is nil")
		assert.NotEmpty(product, "ID does not exist and the object is empty.") // ID does not exist

		assert.Equal(testID, product.ID, "Product ID and test ID do not match.")
	})

}

func TestUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		invalidProduct := models.Product{ID: invalidID}
		assert.False(IsValidUUID(invalidID))
		err := ProductRepo.UpdateProduct(invalidProduct)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "e323a287-c350-4b27-a567-d8c92c52f1d9"
		randomProduct := models.Product{ID: randomUUID, IDC: utils.TestProduct.IDC, Name: utils.TestProduct.Name, Price: utils.TestProduct.Price, Quantity: utils.TestProduct.Quantity}
		assert.True(IsValidUUID(randomUUID))
		err := ProductRepo.UpdateProduct(randomProduct)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		utils.TestProduct.Name = "UPDATED Name"

		err := ProductRepo.UpdateProduct(utils.TestProduct)

		assert.NoError(err, "Product was not updated.")
	})

}

func TestDeleteProduct(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		err := ProductRepo.DeleteProduct(invalidID)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "7d91a563-3386-4069-b785-09c52b5201b5"
		assert.True(IsValidUUID(randomUUID))
		err := ProductRepo.DeleteProduct(randomUUID)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		err := ProductRepo.DeleteProduct(utils.TestProduct.ID)

		assert.NoError(err, "Product was not deleted.")
	})
}

package repositories

import (
	"internship_project/models"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		DropTables(CompanyRepo.DB)
		defer SetupTables(CompanyRepo.DB)
		err := CompanyRepo.AddCompany(&testCompany)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldCompanies, _ := CompanyRepo.GetAllCompanies()
		err := CompanyRepo.AddCompany(&testCompany)
		newCompanies, _ := CompanyRepo.GetAllCompanies()

		assert.NoError(err)
		assert.Equal(len(newCompanies)-len(oldCompanies), 1, "Company was not added.")

	})

}

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful query", func(t *testing.T) {
		allCompanies, err := CompanyRepo.GetAllCompanies()

		assert.NoError(err, "Error was thrown while reading companies")
		assert.NotNil(allCompanies, "Companies returned are nil.")
		assert.IsType(allCompanies, []models.Company{}, "Returned result is not of type Company")

	})
}

func TestGetCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		DropTables(CompanyRepo.DB)
		defer SetupTables(CompanyRepo.DB)
		_, err := CompanyRepo.GetCompany(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while getting from non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		_, err := CompanyRepo.GetCompany(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		_, err := CompanyRepo.GetCompany(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		CompanyRepo.AddCompany(&testCompany)
		company, err := CompanyRepo.GetCompany(testCompany.Id)
		assert.NotNil(company, "Result is nil")
		assert.NoError(err, "There was error while getting company")
		assert.Equal(testCompany.Id, company.Id, "Returned company ID and test ID do not match.")
	})
}

func TestUpdateCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		DropTables(CompanyRepo.DB)
		defer SetupTables(CompanyRepo.DB)
		err := CompanyRepo.UpdateCompany(testCompany)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		testCompany.Id = uuid
		err := CompanyRepo.UpdateCompany(testCompany)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		testCompany.Id = uuid
		err := CompanyRepo.UpdateCompany(testCompany)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		CompanyRepo.AddCompany(&testCompany)
		testCompany.Name = "Updated name"
		err := CompanyRepo.UpdateCompany(testCompany)
		assert.NoError(err, "Company was not updated.")
	})

}

func TestDeleteCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		DropTables(CompanyRepo.DB)
		defer SetupTables(CompanyRepo.DB)
		err := CompanyRepo.DeleteCompany(uuid.NewV4().String())
		assert.Error(err, "Error was not thrown while deleting in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		err := CompanyRepo.DeleteCompany(uuid)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		err := CompanyRepo.DeleteCompany(uuid)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		CompanyRepo.AddCompany(&testCompany)
		err := CompanyRepo.DeleteCompany(testCompany.Id)
		assert.NoError(err, "Company was not deleted.")
	})
}

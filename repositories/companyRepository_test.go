package repositories

import (
	"internship_project/models"
	"internship_project/utils"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := CompanyRepo.AddCompany(&utils.TestCompany)
		assert.Error(err, "Error was not thrown while inserting in non-existing table")
	})

	t.Run("successful query", func(t *testing.T) {
		oldCompanies, _ := CompanyRepo.GetAllCompanies()
		err := CompanyRepo.AddCompany(&utils.TestCompany)
		newCompanies, _ := CompanyRepo.GetAllCompanies()

		assert.NoError(err)
		assert.Equal(1, len(newCompanies)-len(oldCompanies), "Company was not added.")

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
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
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
		CompanyRepo.AddCompany(&utils.TestCompany)
		company, err := CompanyRepo.GetCompany(utils.TestCompany.ID)
		assert.NotNil(company, "Result is nil")
		assert.NoError(err, "There was error while getting company")
		assert.Equal(utils.TestCompany.ID, company.ID, "Returned company ID and test ID do not match.")
	})
}

func TestUpdateCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := CompanyRepo.UpdateCompany(utils.TestCompany)
		assert.Error(err, "Error was not thrown while updating in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		utils.TestCompany.ID = uuid
		err := CompanyRepo.UpdateCompany(utils.TestCompany)
		assert.NotNil(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		utils.TestCompany.ID = uuid
		err := CompanyRepo.UpdateCompany(utils.TestCompany)
		assert.NotNil(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		CompanyRepo.AddCompany(&utils.TestCompany)
		utils.TestCompany.Name = "Updated name"
		err := CompanyRepo.UpdateCompany(utils.TestCompany)
		assert.NoError(err, "Company was not updated.")
	})

}

func TestDeleteCompany(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
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
		CompanyRepo.AddCompany(&utils.TestCompany)
		err := CompanyRepo.DeleteCompany(utils.TestCompany.ID)
		assert.NoError(err, "Company was not deleted.")
	})
}

func TestChangeExternalRightApproveStatus(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(Connpool)
		defer utils.SetUpTables(Connpool)
		err := CompanyRepo.ChangeExternalRightApproveStatus(uuid.NewV4().String(), true)
		assert.Error(err, "Error was not thrown while updating ear in non-existing table")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		uuid := "invalidUUID"
		err := CompanyRepo.ChangeExternalRightApproveStatus(uuid, true)
		assert.Error(err, "Error was not thrown for invalid uuid")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		err := CompanyRepo.ChangeExternalRightApproveStatus(uuid, true)
		assert.Error(err, "Error was not thrown for non-existing uuid")
	})

	t.Run("successful query", func(t *testing.T) {
		err := CompanyRepo.ChangeExternalRightApproveStatus(utils.TestEar1.ID, true)
		assert.NoError(err, "Ear was not approved.")
	})
}

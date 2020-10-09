package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCompany(t *testing.T) {
	assert := assert.New(t)

	assert.True(DoesTableExist("companies", CompanyRepo.DB), "Table does not exist.")

	oldCompanies, _ := CompanyRepo.GetAllCompanies()
	err := CompanyRepo.AddCompany(&testCompany)
	newCompanies, _ := CompanyRepo.GetAllCompanies()

	assert.NoError(err)
	assert.Equal(len(newCompanies)-len(oldCompanies), 1, "Company was not added.")
}

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)
	allCompanies, err := CompanyRepo.GetAllCompanies()

	assert.NoError(err)
	assert.NotNil(allCompanies, "Companies returned were nil.")
	assert.IsType(allCompanies, []models.Company{})
}

func TestGetCompany(t *testing.T) {
	assert := assert.New(t)
	testID := testCompany.Id

	assert.True(IsValidUUID(testID), "Company ID is not valid.")

	company, err := CompanyRepo.GetCompany(testID)

	assert.NoError(err)
	assert.NotNil(company, "Company is nil.")
	assert.NotEmpty(company, "Company ID does not exist.")

	assert.Equal(testID, company.Id, "Company ID and test ID do not match.")
}

func TestUpdateCompany(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testCompany.Id), "Company ID is not valid.")

	companyForUpdate, _ := CompanyRepo.GetCompany(testCompany.Id)
	companyForUpdate.Name = "UPDATED NAME"

	err := CompanyRepo.UpdateCompany(companyForUpdate)

	assert.NoError(err, "Company was not updated.")
}

func TestDeleteCompany(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testCompany.Id), "Company ID is not valid.")

	err := CompanyRepo.DeleteCompany(testCompany.Id)

	assert.NoError(err, "Company was not deleted.")
}

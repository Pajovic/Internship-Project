package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCompany(t *testing.T) {
	assert := assert.New(t)

	oldCompanies, _ := repository.GetAllCompanies()
	err := repository.AddCompany(&testCompany)
	newCompanies, _ := repository.GetAllCompanies()

	assert.NoError(err)
	assert.Equal(len(newCompanies)-len(oldCompanies), 1)
}

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)
	allCompanies, err := repository.GetAllCompanies()

	assert.NoError(err)
	assert.NotNil(allCompanies)
	assert.IsType(allCompanies, []models.Company{})
}

func TestGetCompany(t *testing.T) {
	assert := assert.New(t)
	testID := testCompany.Id

	comapny, err := repository.GetCompany(testID)

	assert.NoError(err)
	assert.NotNil(comapny)
	assert.NotEmpty(comapny)

	assert.Equal(testID, comapny.Id)
}

func TestUpdateCompany(t *testing.T) {
	assert := assert.New(t)

	companyForUpdate, _ := repository.GetCompany(testCompany.Id)
	companyForUpdate.Name = "UPDATED NAME"

	err := repository.UpdateCompany(companyForUpdate)

	assert.NoError(err)
}

func TestDeleteCompany(t *testing.T) {
	assert := assert.New(t)

	err := repository.DeleteCompany(testCompany.Id)

	assert.NoError(err)
}

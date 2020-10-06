package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)
	allCompanies, err := GetAllCompanies(connection)

	assert.NoError(err)
	assert.NotNil(allCompanies)
	assert.IsType(allCompanies, []models.Company{})
}

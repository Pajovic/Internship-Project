package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	newEmployee := models.Employee{
		FirstName: "Test Name",
		LastName:  "Test Surname",
		CompanyID: "153fac6d-760d-4841-87e9-15aee2f25182", // ID from database
		C:         false,
		R:         true,
		U:         false,
		D:         false,
	}

	oldEmployees, _ := GetAllEmployees(DbInstance)
	err := AddEmployee(DbInstance, &newEmployee)
	newEmployees, _ := GetAllEmployees(DbInstance)

	assert.NoError(err)
	assert.Equal(len(newEmployees)-len(oldEmployees), 1)
}

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)
	allEmployees, err := GetAllEmployees(DbInstance)

	assert.NoError(err)
	assert.NotNil(allEmployees)
	assert.IsType(allEmployees, []models.Employee{})
}

func TestGetEmployeeByID(t *testing.T) {
	assert := assert.New(t)

	allEmployees, _ := GetAllEmployees(DbInstance)
	if len(allEmployees) == 0 {
		require.Fail(t, "There are no employees to get.")
	}
	employeeID := allEmployees[0].ID

	employee, err := GetEmployeeByID(DbInstance, employeeID)

	assert.NoError(err)
	assert.NotNil(employee)
	assert.NotEmpty(employee)

	assert.Equal(employeeID, employee.ID)
}

func TestUpdateEmployee(t *testing.T) {
	assert := assert.New(t)

	allEmployees, _ := GetAllEmployees(DbInstance)
	if len(allEmployees) == 0 {
		require.Fail(t, "There are no employees for update.")
	}
	employeeForUpdate := allEmployees[0]
	employeeForUpdate.LastName = "UPDATED Last Name"

	err := UpdateEmployee(DbInstance, employeeForUpdate)

	assert.NoError(err)
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	allEmployees, _ := GetAllEmployees(DbInstance)
	if len(allEmployees) == 0 {
		require.Fail(t, "There are no employees for deletion.")
	}
	employeeForDeletion := allEmployees[0]

	err := DeleteEmployee(DbInstance, employeeForDeletion.ID)

	assert.NoError(err)
}

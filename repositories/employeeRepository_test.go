package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	oldEmployees, _ := GetAllEmployees(DbInstance)
	err := AddEmployee(DbInstance, &testEmployee)
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
	testID := testEmployee.ID

	employee, err := GetEmployeeByID(DbInstance, testID)

	assert.NoError(err)
	assert.NotNil(employee)
	assert.NotEmpty(employee)

	assert.Equal(testID, employee.ID)
}

func TestUpdateEmployee(t *testing.T) {
	assert := assert.New(t)

	employeeForUpdate, _ := GetEmployeeByID(DbInstance, testEmployee.ID)
	employeeForUpdate.LastName = "UPDATED Last Name"

	err := UpdateEmployee(DbInstance, employeeForUpdate)

	assert.NoError(err)
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	err := DeleteEmployee(DbInstance, testEmployee.ID)

	assert.NoError(err)
}

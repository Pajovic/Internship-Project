package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	oldEmployees, _ := EmployeeRepo.GetAllEmployees()
	err := EmployeeRepo.AddEmployee(&testEmployee)
	newEmployees, _ := EmployeeRepo.GetAllEmployees()

	assert.NoError(err)
	assert.Equal(len(newEmployees)-len(oldEmployees), 1)
}

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)
	allEmployees, err := EmployeeRepo.GetAllEmployees()

	assert.NoError(err)
	assert.NotNil(allEmployees)
	assert.IsType(allEmployees, []models.Employee{})
}

func TestGetEmployeeByID(t *testing.T) {
	assert := assert.New(t)
	testID := testEmployee.ID

	employee, err := EmployeeRepo.GetEmployeeByID(testID)

	assert.NoError(err)
	assert.NotNil(employee)
	assert.NotEmpty(employee)

	assert.Equal(testID, employee.ID)
}

func TestUpdateEmployee(t *testing.T) {
	assert := assert.New(t)

	employeeForUpdate, _ := EmployeeRepo.GetEmployeeByID(testEmployee.ID)
	employeeForUpdate.LastName = "UPDATED Last Name"

	err := EmployeeRepo.UpdateEmployee(employeeForUpdate)

	assert.NoError(err)
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	err := EmployeeRepo.DeleteEmployee(testEmployee.ID)

	assert.NoError(err)
}

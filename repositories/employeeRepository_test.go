package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	assert.True(DoesTableExist("employees", EmployeeRepo.DB), "Table does not exist.")

	oldEmployees, _ := EmployeeRepo.GetAllEmployees()
	err := EmployeeRepo.AddEmployee(&testEmployee)
	newEmployees, _ := EmployeeRepo.GetAllEmployees()

	assert.NoError(err)
	assert.Equal(len(newEmployees)-len(oldEmployees), 1, "Employee was not added.")
}

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)
	allEmployees, err := EmployeeRepo.GetAllEmployees()

	assert.NoError(err)
	assert.NotNil(allEmployees, "Employees returned were nil.")
	assert.IsType(allEmployees, []models.Employee{})
}

func TestGetEmployeeByID(t *testing.T) {
	assert := assert.New(t)
	testID := testEmployee.ID

	assert.True(IsValidUUID(testID), "Employee ID is not valid.")

	employee, err := EmployeeRepo.GetEmployeeByID(testID)

	assert.NoError(err)
	assert.NotNil(employee, "Employee returned was nil.")
	assert.NotEmpty(employee, "Employee ID does not exist.") // ID does not exist

	assert.Equal(testID, employee.ID, "Employee ID and test ID do not match.")
}

func TestUpdateEmployee(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testEmployee.ID), "Employee ID is not valid.")

	employeeForUpdate, _ := EmployeeRepo.GetEmployeeByID(testEmployee.ID)
	employeeForUpdate.LastName = "UPDATED Last Name"

	err := EmployeeRepo.UpdateEmployee(employeeForUpdate)

	assert.NoError(err, "Employee was not updated.")
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidUUID(testEmployee.ID), "Employee ID is not valid.")

	err := EmployeeRepo.DeleteEmployee(testEmployee.ID)

	assert.NoError(err, "Employee was not deleted.")
}

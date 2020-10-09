package repositories

import (
	"internship_project/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		DropTables(EmployeeRepo.DB)
		defer SetupTables(EmployeeRepo.DB)
		assert.False(DoesTableExist("employees", EmployeeRepo.DB))
	})

	t.Run("table was added back", func(t *testing.T) {
		SetupTables(EmployeeRepo.DB)
		assert.True(DoesTableExist("employees", EmployeeRepo.DB))
	})

	t.Run("successful query", func(t *testing.T) {
		oldEmployees, _ := EmployeeRepo.GetAllEmployees()
		err := EmployeeRepo.AddEmployee(&testEmployee)
		newEmployees, _ := EmployeeRepo.GetAllEmployees()

		assert.NoError(err)
		assert.Equal(len(newEmployees)-len(oldEmployees), 1, "Employee was not added.")
	})
}

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful GetAll query", func(t *testing.T) {
		allEmployees, err := EmployeeRepo.GetAllEmployees()
		assert.NoError(err)
		assert.NotNil(allEmployees, "Employees returned were nil.")
		assert.IsType(allEmployees, []models.Employee{})
	})

}

func TestGetEmployeeByID(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		_, err := EmployeeRepo.GetEmployeeByID(invalidID)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "c5ef08c6-60eb-4687-bcbb-df37ebc9e105"
		assert.True(IsValidUUID(randomUUID))
		_, err := EmployeeRepo.GetEmployeeByID(randomUUID)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		testID := testEmployee.ID
		employee, err := EmployeeRepo.GetEmployeeByID(testID)

		assert.NoError(err)
		assert.NotNil(employee, "Employee returned was nil.")
		assert.NotEmpty(employee, "Employee ID does not exist.") // ID does not exist

		assert.Equal(testID, employee.ID, "Employee ID and test ID do not match.")
	})
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

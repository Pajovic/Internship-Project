package repositories

import (
	"internship_project/models"
	"internship_project/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(EmployeeRepo.DB)
		defer SetupTables(EmployeeRepo.DB)

		assert.False(DoesTableExist("employees", EmployeeRepo.DB))
		err := EmployeeRepo.AddEmployee(&testEmployee)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		oldEmployees, _ := EmployeeRepo.GetAllEmployees(testAdmin.CompanyID)
		err := EmployeeRepo.AddEmployee(&testEmployee)
		newEmployees, _ := EmployeeRepo.GetAllEmployees(testAdmin.CompanyID)

		assert.NoError(err)
		assert.Equal(1, len(newEmployees)-len(oldEmployees), "Employee was not added.")
	})

	t.Run("add an existing employee", func(t *testing.T) {
		existingEmployee := &models.Employee{ID: testEmployee.ID}
		err := EmployeeRepo.AddEmployee(existingEmployee)

		assert.Error(err)
	})
}

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)

	t.Run("successful GetAll query", func(t *testing.T) {
		allEmployees, err := EmployeeRepo.GetAllEmployees(testAdmin.ID)

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

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		invalidEmployee := models.Employee{ID: invalidID, FirstName: "Test", LastName: "Test", CompanyID: testEmployee.CompanyID}
		assert.False(IsValidUUID(invalidID))
		err := EmployeeRepo.UpdateEmployee(invalidEmployee)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "7d91a563-3386-4069-b785-09c52b5201b5"
		randomEmployee := models.Employee{ID: randomUUID, FirstName: "Test", LastName: "Test", CompanyID: testEmployee.CompanyID}
		assert.True(IsValidUUID(randomUUID))
		err := EmployeeRepo.UpdateEmployee(randomEmployee)
		assert.Error(err)
		assert.Equal("No row found to update", err.Error())
	})

	t.Run("successful query", func(t *testing.T) {
		employeeForUpdate, _ := EmployeeRepo.GetEmployeeByID(testEmployee.ID)
		employeeForUpdate.LastName = "UPDATED Last Name"

		err := EmployeeRepo.UpdateEmployee(employeeForUpdate)

		assert.NoError(err, "Employee was not updated.")
	})
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		err := EmployeeRepo.DeleteEmployee(invalidID)
		assert.Error(err)
	})

	t.Run("id does not exist", func(t *testing.T) {
		randomUUID := "7d91a563-3386-4069-b785-09c52b5201b5"
		assert.True(IsValidUUID(randomUUID))
		err := EmployeeRepo.DeleteEmployee(randomUUID)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		err := EmployeeRepo.DeleteEmployee(testEmployee.ID)

		assert.NoError(err, "Employee was not deleted.")
	})
}

func TestGetEmployeeExternalPermissions(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		_, err := EmployeeRepo.GetEmployeeExternalPermissions(invalidID, testProduct)
		assert.Error(err)
	})

	t.Run("company with id does not exist", func(t *testing.T) {
		randomUUID := "7d91a563-3386-4069-b785-09c52b5201b5"
		assert.True(IsValidUUID(randomUUID))
		_, err := EmployeeRepo.GetEmployeeExternalPermissions(randomUUID, testProduct)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		_, err := EmployeeRepo.GetEmployeeExternalPermissions(testCompany2.Id, testProduct)

		assert.NoError(err, "Could not get employee's EAR")
	})
}

func TestCheckCompaniesSharingEmployeeData(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid id", func(t *testing.T) {
		invalidID := "123-asd-321"
		assert.False(IsValidUUID(invalidID))
		err := EmployeeRepo.CheckCompaniesSharingEmployeeData(invalidID, testProduct.IDC)
		assert.Error(err)
	})

	t.Run("company with id does not exist", func(t *testing.T) {
		randomUUID := "7d91a563-3386-4069-b785-09c52b5201b5"
		assert.True(IsValidUUID(randomUUID))
		err := EmployeeRepo.CheckCompaniesSharingEmployeeData(randomUUID, testProduct.IDC)
		assert.Error(err)
	})

	t.Run("successful query", func(t *testing.T) {
		err := EmployeeRepo.CheckCompaniesSharingEmployeeData(testCompany2.Id, testCompany1.Id)

		assert.NoError(err, "Companies do not share data")
	})
}

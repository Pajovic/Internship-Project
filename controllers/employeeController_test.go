package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internship_project/models"
	"internship_project/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllEmployees(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/employee", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("employeeID", utils.TestAdmin.ID)

	handler := http.HandlerFunc(EmployeeCont.GetAllEmployees)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, employees, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestGetEmployeeById(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", EmployeeCont.GetEmployeeByID)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		path := fmt.Sprintf("/employee/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("employeeID", utils.TestAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/employee/%s", "INVALID_UUID")
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("employeeID", utils.TestAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/employee/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("employeeID", utils.TestAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		path := fmt.Sprintf("/employee/%s", utils.TestAdmin.ID)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("employeeID", utils.TestAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestDeleteEmployee(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", EmployeeCont.DeleteEmployee)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		path := fmt.Sprintf("/employee/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.employees, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/employee/%s", "INVALID_UUID")
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/employee/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful delete", func(t *testing.T) {
		defer utils.SetUpTables(connpool)

		path := fmt.Sprintf("/employee/%s", utils.TestEmployee1.ID)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNoContent, rr.Code, "Response code is not correct")
	})
}

func TestAddEmployee(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(EmployeeCont.AddNewEmployee)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEmployee)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.employees, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful add", func(t *testing.T) {
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEmployee)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestUpdateEmployee(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/employee", EmployeeCont.UpdateEmployee)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestAdmin)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/employee")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.employees, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		invalidUUID := models.Employee{
			ID: "INVALID_UUID",
		}
		body, err := json.Marshal(invalidUUID)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/employee")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		invalidUUID := models.Employee{
			ID: uuid.NewV4().String(),
		}
		body, err := json.Marshal(invalidUUID)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/employee")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful update", func(t *testing.T) {
		utils.TestAdmin.FirstName = "UPDATED"

		body, err := json.Marshal(utils.TestAdmin)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/employee")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

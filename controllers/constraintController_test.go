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

func TestGetAllConstraints(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/constraint", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(ConstraintCont.GetAllConstraints)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.access_constraints, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestGetConstraintById(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/constraint/{id}", ConstraintCont.GetConstraintById)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/constraint/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, access_constraints, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/constraint/%s", "INVALID_UUID")
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/constraint/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		path := fmt.Sprintf("/constraint/%s", testConstraint2.ID)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestUpdateConstraint(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/constraint", ConstraintCont.UpdateConstraint)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		body, err := json.Marshal(testConstraint)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/constraint")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.access_constraints, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		invalidUUID := models.AccessConstraint{
			ID: "INVALID_UUID",
		}
		body, err := json.Marshal(invalidUUID)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/constraint")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
		t.Log(rr.Body.String())
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		randomUUID := models.AccessConstraint{
			ID:            uuid.NewV4().String(),
			IDEAR:         testConstraint2.IDEAR,
			OperatorID:    testConstraint2.OperatorID,
			PropertyID:    testConstraint2.PropertyID,
			PropertyValue: testConstraint2.PropertyValue,
		}
		body, err := json.Marshal(randomUUID)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/constraint")
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
		defer SetUpTables(connpool)

		testConstraint2.PropertyValue = 99
		body, err := json.Marshal(testConstraint2)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/constraint")
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

func TestDeleteConstraint(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/constraint/{id}", ConstraintCont.DeleteConstraint)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/constraint/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.access_constraints, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/constraint/%s", "INVALID_UUID")
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/constraint/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful delete", func(t *testing.T) {
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/constraint/%s", testConstraint2.ID)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(204, rr.Code, "Response code is not correct")
	})
}

func TestAddConstraint(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(ConstraintCont.AddConstraint)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		body, err := json.Marshal(testConstraint)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/constraint", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.access_constraints, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful add", func(t *testing.T) {
		defer SetUpTables(connpool)

		body, err := json.Marshal(testConstraint)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/constraint", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

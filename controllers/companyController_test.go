package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internship_project/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCompanies(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/company", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(CompanyCont.GetAllCompanies)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.companies, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestGetCompanyById(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/company/{id}", CompanyCont.GetCompanyById)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/company/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.companies, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/%s", "INVALID_UUID")
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		path := fmt.Sprintf("/company/%s", testCompany1.Id)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestDeleteCompany(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/company/{id}", CompanyCont.DeleteCompany)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/company/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.companies, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/%s", "INVALID_UUID")
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful delete", func(t *testing.T) {
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/company/%s", mainCompany1.Id)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(204, rr.Code, "Response code is not correct")
	})
}

func TestAddCompany(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(CompanyCont.AddCompany)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		body, err := json.Marshal(testCompany)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/company", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.companies, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful add", func(t *testing.T) {
		defer SetUpTables(connpool)

		body, err := json.Marshal(testCompany)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/company", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestUpdateCompany(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/company", CompanyCont.UpdateCompany)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		body, err := json.Marshal(testCompany1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/company")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.companies, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		body, err := json.Marshal(testCompany)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/company")
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
		testCompany.Id = uuid.NewV4().String()
		body, err := json.Marshal(testCompany)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/company")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		testCompany.Id = ""

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful update", func(t *testing.T) {
		defer SetUpTables(connpool)

		body, err := json.Marshal(testCompany1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/company")
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

func TestChangeExternalRightApproveStatus(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/company/approve/{idear}", func(w http.ResponseWriter, r *http.Request) {
		CompanyCont.ChangeExternalRightApproveStatus(w, r, true)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/approve/%s", "INVALID_UUID")
		req, err := http.NewRequest("PATCH", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("companyID", mainCompany1.Id)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusBadRequest, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/company/approve/%s", uuid.NewV4().String())
		req, err := http.NewRequest("PATCH", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("companyID", mainCompany1.Id)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful update", func(t *testing.T) {
		path := fmt.Sprintf("/company/approve/%s", testEar2.ID)
		req, err := http.NewRequest("PATCH", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("companyID", mainCompany1.Id)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})

	t.Run("no permission to update", func(t *testing.T) {
		path := fmt.Sprintf("/company/approve/%s", testEar2.ID)
		req, err := http.NewRequest("PATCH", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("companyID", testCompany1.Id)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
		assert.Equal(`Your company does not have permission to approve sharing`, rr.Body.String(), "Error message is not correct")
	})
}

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

func TestGetAllEars(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/ear", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(ExternalRightCont.GetAllEars)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.external_access_rights, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestGetExternalRightById(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/ear/{id}", ExternalRightCont.GetEarById)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		path := fmt.Sprintf("/ear/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, external_access_rights, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/ear/%s", "INVALID_UUID")
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/ear/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		path := fmt.Sprintf("/ear/%s", utils.TestEar1.ID)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestUpdateExternalRight(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/ear", ExternalRightCont.UpdateEar)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEar1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/ear")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.external_access_rights, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		invalidUUID := models.ExternalRights{
			ID:     "INVALID_UUID",
			Read:   true,
			Update: true,
			Delete: true,
			IDSC:   utils.TestCompany1.ID,
			IDRC:   utils.TestCompany2.ID,
		}
		body, err := json.Marshal(invalidUUID)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/ear")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(500, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		body, err := json.Marshal(utils.TestEar)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/ear")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusNotFound, rr.Code, "Response code is not correct")
	})

	t.Run("successful update", func(t *testing.T) {
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEar1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/ear")
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

func TestDeleteExternalRight(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/ear/{id}", ExternalRightCont.DeleteEar)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		path := fmt.Sprintf("/ear/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.external_access_rights, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/ear/%s", "INVALID_UUID")
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/ear/%s", uuid.NewV4().String())
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

		path := fmt.Sprintf("/ear/%s", utils.TestEar2.ID)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(204, rr.Code, "Response code is not correct")
	})
}

func TestAddExternalRight(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(ExternalRightCont.AddEar)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEar)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/ear", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, public.external_access_rights, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("successful add", func(t *testing.T) {
		defer utils.SetUpTables(connpool)

		body, err := json.Marshal(utils.TestEar)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/ear", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

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

func TestGetAllProducts(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/product", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("employeeID", testAdmin.ID)

	handler := http.HandlerFunc(ProductCont.GetAllProducts)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestAddProduct(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(ProductCont.AddProduct)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		body, err := json.Marshal(testProduct)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful add", func(t *testing.T) {
		defer SetUpTables(connpool)

		body, err := json.Marshal(testProduct)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestGetProductById(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/product/{id}", ProductCont.GetProductById)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/product/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
		assert.Equal(`The table you wish to work with, products, does not exist.`, rr.Body.String(), "Error message is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/product/%s", "INVALID_UUID")
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/product/%s", uuid.NewV4().String())
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful get", func(t *testing.T) {
		path := fmt.Sprintf("/product/%s", testProduct1.ID)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestDeleteProduct(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/product/{id}", ProductCont.DeleteProduct)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/product/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		path := fmt.Sprintf("/product/%s", "INVALID_UUID")
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		path := fmt.Sprintf("/product/%s", uuid.NewV4().String())
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful delete", func(t *testing.T) {
		defer SetUpTables(connpool)

		path := fmt.Sprintf("/product/%s", testProduct1.ID)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

func TestUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	router := mux.NewRouter()
	router.HandleFunc("/product", ProductCont.UpdateProduct)

	t.Run("table does not exist", func(t *testing.T) {
		utils.DropTables(connpool)
		defer SetUpTables(connpool)

		testProduct1.Name = "UPDATED"

		body, err := json.Marshal(testProduct1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/product")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		body, err := json.Marshal(testProduct)
		if err != nil {
			t.Fatal(err)
		}

		path := fmt.Sprintf("/product")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("non-existing uuid", func(t *testing.T) {
		testProduct.ID = uuid.NewV4().String()
		body, err := json.Marshal(testProduct)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/product")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		testProduct.ID = ""

		assert.Equal(http.StatusInternalServerError, rr.Code, "Response code is not correct")
	})

	t.Run("successful update", func(t *testing.T) {
		defer SetUpTables(connpool)

		testProduct1.Name = "UPDATED"

		body, err := json.Marshal(testProduct1)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/product")
		req, err := http.NewRequest("PUT", path, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("employeeID", testAdmin.ID)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Code, "Response code is not correct")
	})
}

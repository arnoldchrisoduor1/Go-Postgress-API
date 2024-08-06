// When running these tests we need to ensure that the database is properly set-up.
package main_test

import (
	"testing"
	"os"
	"log"
	"github.com/arnoldchrisoduor1/Go-Postgress-API"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
)

var a main.App

func TestMain(m *testing.M) {
	a.Initialize (
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

// This function checks if the function exists.
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// This function deletes everything from the table.
func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}


// Testing the response to the /products endpoint with an empty table.
func TestEmptyTable(t *testing.T) {
	clearTable()	// deletes all records from the products table.

	req, _:= http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	// Now we test if the HTTP response is what we expect it to be
	checkResponseCode(t, http.StatusOK, response.Code)


	// Checking the body of the response and test it is the textual representation of an empty array.
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Function to check for non-existent products.
func TestGetNonExistentProduct(t *testing.T) {
	clearTable()	// Ensuring the table is empty.

	// Tries to access a non-existent product at an endpoint.
	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	// status code is 404, indicating the product was not found.
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

// This test code will create a product.
func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name":"test product", "price":11.22}`)

	// We manually add an item to the databse.
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Accesssing the relevant endpoints to fetch the product.
	response := executeRequest(req)
	// Checking for status code 201 for resource created.
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&m)

	// Checking if the response contained a JSON object with contents identical to that of the payload.

	if m["name"] != "test product" {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is converted to 1.0 because JSON unmarshaling coneverts numbers to floats,
	// when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

// Function to execute our request.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// Implementing the checkResponseCode function.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)`
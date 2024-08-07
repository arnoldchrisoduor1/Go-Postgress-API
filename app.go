// a simple application to test the connections to the databases.

package main

// importing the necessary packages.
import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Creating the struct App.
type App struct {
	Router *mux.Router	// pointer to 'mux.Router' to handle HTTP routing.
	DB     *sql.DB		// pointer to sql.DB which represents a connection to the databse.
}

// A method on the 'App' struct that initializes the application.
func (a *App) Initialize() { 
	connectionString :=
		"host=localhost port=5432 user=postgres password=12345 dbname=go_api sslmode=disable"

	var err error // declares variable err of the type error.
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err) // Returns an error if there's a problem logging into the database.
	}
	a.Router = mux.NewRouter() // Initializes a new Gorilla mux router and assigns it to 'a.Router'

	a.InitializeRoutes()
}

func (a *App) Run(addr string) {	// Defines a method on the 'App' struct to start the HTTP server.
	log.Fatal(http.ListenAndServe(":8010", a.Router)) // Logs fatal error if the server fails to start.
}

// Creating a handler for the route that fetches a single product.
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Extracting variables from the request's URL path.
	id, err := strconv.Atoi(vars["id"]) // converting the extracted "id" to string.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id} // initializes a 'product' struct with the 'ID' field set to the extracted 'id'.
	if err := p.getProduct(a.DB); err != nil {	// calls the 'getProduct' method on the 'product' from the database.
		switch err {
		case sql.ErrNoRows:
			// If products are not found.
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			// Respond with internal server error if there's any other error.
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	responseWithJSON(w, http.StatusOK, p)	// responds with the product in JSON format if no error occurs.
}


// A helper function to respond with an error message in JSON format.
func respondWithError(w http.ResponseWriter, code int, message string) {
	responseWithJSON(w, code, map[string]string{"error":message})
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload) // Marshalls the payload into JSON format.

	w.Header().Set("Content-Type", "application/json") // sets the response content-type to JSON.
	w.WriteHeader(code)		// sets the HTTP status code.
	w.Write(response)	// Writes the JSON response.
}

// A handler to fetch a list of products.
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	
	count, _ := strconv.Atoi(r.FormValue("count")) // converts count form value to interger, defaults to 0 if fails.
	start, _ := strconv.Atoi(r.FormValue("start")) // converts start form value to an interger, defaults to 10 if fail.

	if count > 10 || count  < 1 {
		// Default to 10 if it's outside range of 1 to 10
		count = 10
	}
	if start < 0 {
		// sets default start value to 0 if it's negative.
		start = 0
	}

	products, err := getProducts(a.DB, start, count) // Calls the getProducts function to get the products from DB.
	if err != nil {
		// return internal server error if query fails.
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, products) // JSON response if no error occurs.
}

// A handler to create a product.
func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {

	var p product

	// Assumes the request body is a json object containing the details of the product to be created.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJSON(w, http.StatusCreated, p)
}

// A handler to update a product.
func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, p)
}

// A handler to delete a product from the database.
func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, map[string]string{"result":"success"})
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}
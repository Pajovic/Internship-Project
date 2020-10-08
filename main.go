package main

import (
	"context"
	"fmt"
	"internship_project/controllers"
	"internship_project/repositories"
	"internship_project/services"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
)

type config struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DatabaseURL string `json:"database_url"`
}

func main() {
	connpool := getConnectionPool()
	productController := getProductController(connpool)
	defer connpool.Close()

	r := mux.NewRouter()

	productRouter := r.PathPrefix("/product").Subrouter()

	productRouter.HandleFunc("", productController.GetAllProducts).Methods("GET")

	productRouter.HandleFunc("/{id}", productController.GetProductById).Methods("GET")

	productRouter.HandleFunc("", productController.AddProduct).Methods("POST")

	productRouter.HandleFunc("/{id}", productController.UpdateProduct).Methods("PUT")

	productRouter.HandleFunc("/{id}", productController.DeleteProduct).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

func getConnectionPool() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("database.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.DatabaseURL)

	connection, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to database.")

	return connection
}

func getProductController(connpool *pgxpool.Pool) controllers.ProductController {

	productRepository := repositories.ProductRepository{DB: connpool}
	productService := services.ProductService{Repository: productRepository}
	productController := controllers.ProductController{Service: productService}

	fmt.Println("Employee controller up and running.")

	return productController
}

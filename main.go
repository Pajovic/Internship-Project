package main

import (
	"context"
	"fmt"
	"github.com/codingsince1985/geo-golang/mapquest/nominatim"
	"internship_project/controllers"
	"internship_project/elasticsearch_helpers"
	"internship_project/kafka_helpers"
	"internship_project/repositories"
	"internship_project/services"
	"internship_project/utils"
	"net/http"
	"os"

	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
	"github.com/segmentio/kafka-go"
)

type DbConfig struct {
	DbUsername      string `json:"db_username"`
	DbPassword      string `json:"db_password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}


type NominatimConfig struct {
	Key string `json:"nominatim_key"`
}

type KafkaEsConfig struct {
	MainKafkaTopic  string `json:"main_kafka_topic"`
	RetryKafkaTopic string `json:"retry_kafka_topic"`
	MainTopicTime   int    `json:"main_topic_time"`
	RetryTopicTime  int    `json:"retry_topic_time"`
	KafkaAddress    string `json:"kafka_address"`
	KafkaGroupId    string `json:"kafka_group_id"`
	EsAddress       string `json:"es_address"`
}

var (
	userRepository repositories.UserRepository
	userService    services.UserService
)

func main() {
	var db_conf DbConfig
	if _, err := confl.DecodeFile("dbconfig.conf", &db_conf); err != nil {
		panic(err)
	}

	var kafka_es_conf KafkaEsConfig
	if _, err := confl.DecodeFile("kafka_es_cofig.conf", &kafka_es_conf); err != nil {
		panic(err)
	}

	connpool := getConnectionPool(db_conf)
	defer connpool.Close()

	kafkaWriter := kafka_helpers.GetWriter("ava-internship")
	defer kafkaWriter.Close()

	EsClient := elasticsearch_helpers.GetElasticsearchClient(kafka_es_conf.EsAddress)
	kafkaConsumer := kafka_helpers.NewConsumer(kafka_es_conf.MainKafkaTopic, kafka_es_conf.KafkaAddress, kafka_es_conf.KafkaGroupId, EsClient, kafka_es_conf.MainTopicTime)
	go kafkaConsumer.Consume()
	defer kafkaConsumer.Reader.Close()

	kafkaRetryHandler := kafka_helpers.GetRetryHandler(kafka_es_conf.RetryKafkaTopic, kafka_es_conf.MainKafkaTopic, kafka_es_conf.KafkaAddress, kafka_es_conf.KafkaGroupId, kafka_es_conf.RetryTopicTime)

	employeeController := getEmployeeController(connpool)
	productController := getProductController(connpool, &employeeController.Service.Repository, kafkaWriter, EsClient)
	companyController := GetCompanyController(connpool, kafkaWriter)
	ExternalRightController := getExternalRightController(connpool)
	constraintController := getConstraintController(connpool)
	userController := getUserController(connpool)
	shopController := getShopController(connpool)

	userRepository = repositories.NewUserRepo(connpool)
	userService = services.UserService{Repository: userRepository}

	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./public/")))
	r.PathPrefix("/static/").Handler(s)

	// Sign In Routes
	r.HandleFunc("/auth/google", userController.GoogleAuth).Methods("POST")

	r.HandleFunc("/search", productController.SearchProducts).Methods("GET")

	// Kafka routes
	kafkaRouter := r.PathPrefix("/kafka").Subrouter()
	kafkaRouter.HandleFunc("/retry", kafkaRetryHandler.TransferToMainTopic).Methods("POST")

	// Product Routes
	productRouter := r.PathPrefix("/product").Subrouter()
	productRouter.Headers("employeeID")

	productRouter.HandleFunc("", productController.GetAllProducts).Methods("GET")
	productRouter.HandleFunc("/{id}", productController.GetProductById).Methods("GET")
	productRouter.HandleFunc("", productController.AddProduct).Methods("POST")
	productRouter.HandleFunc("", productController.UpdateProduct).Methods("PUT")
	productRouter.HandleFunc("/{id}", productController.DeleteProduct).Methods("DELETE")

	// Company Routes
	companyRouter := r.PathPrefix("/company").Subrouter()
	companyRouter.Headers("companyID")

	companyRouter.HandleFunc("", companyController.GetAllCompanies).Methods("GET")
	companyRouter.HandleFunc("/{id}", companyController.GetCompanyById).Methods("GET")
	companyRouter.HandleFunc("", companyController.AddCompany).Methods("POST")
	companyRouter.HandleFunc("", companyController.UpdateCompany).Methods("PUT")
	companyRouter.HandleFunc("/{id}", companyController.DeleteCompany).Methods("DELETE")
	companyRouter.HandleFunc("/approve/{idear}", func(w http.ResponseWriter, r *http.Request) {
		companyController.ChangeExternalRightApproveStatus(w, r, true)
	}).Methods("PATCH")
	companyRouter.HandleFunc("/disapprove/{idear}", func(w http.ResponseWriter, r *http.Request) {
		companyController.ChangeExternalRightApproveStatus(w, r, true)
	}).Methods("PATCH")

	// Employee Routes
	employeeRouter := r.PathPrefix("/employees").Subrouter()
	productRouter.Headers("employeeID")

	employeeRouter.HandleFunc("", employeeController.GetAllEmployees).Methods("GET")
	employeeRouter.HandleFunc("/{id}", employeeController.GetEmployeeByID).Methods("GET")
	employeeRouter.HandleFunc("", employeeController.AddNewEmployee).Methods("POST")
	employeeRouter.HandleFunc("", employeeController.UpdateEmployee).Methods("PUT")
	employeeRouter.HandleFunc("/{id}", employeeController.DeleteEmployee).Methods("DELETE")

	// External Access Rules Routes
	earRouter := r.PathPrefix("/ear").Subrouter()

	earRouter.HandleFunc("", ExternalRightController.GetAllEars).Methods("GET")
	earRouter.HandleFunc("/{id}", ExternalRightController.GetEarById).Methods("GET")
	earRouter.HandleFunc("", ExternalRightController.AddEar).Methods("POST")
	earRouter.HandleFunc("", ExternalRightController.UpdateEar).Methods("PUT")
	earRouter.HandleFunc("/{id}", ExternalRightController.DeleteEar).Methods("DELETE")

	// Constraints Routes
	constraintRouter := r.PathPrefix("/constraint").Subrouter()

	constraintRouter.HandleFunc("", constraintController.GetAllConstraints).Methods("GET")
	constraintRouter.HandleFunc("/{id}", constraintController.GetConstraintById).Methods("GET")
	constraintRouter.HandleFunc("", constraintController.AddConstraint).Methods("POST")
	constraintRouter.HandleFunc("", constraintController.UpdateConstraint).Methods("PUT")
	constraintRouter.HandleFunc("/{id}", constraintController.DeleteConstraint).Methods("DELETE")

	// Shop Routes
	shopRouter := r.PathPrefix("/shop").Subrouter()

	shopRouter.HandleFunc("", shopController.GetAllShops).Methods("GET")
	shopRouter.HandleFunc("/{id}", shopController.GetShopById).Methods("GET")
	shopRouter.HandleFunc("", shopController.AddShop).Methods("POST")
	shopRouter.HandleFunc("", shopController.UpdateShop).Methods("PUT")
	shopRouter.HandleFunc("/{id}", shopController.DeleteShop).Methods("DELETE")
	shopRouter.HandleFunc("/{id}/address", shopController.GetAddress).Methods("GET")

	companyRouter.Use(googleAuthMiddleware)
	constraintRouter.Use(googleAuthMiddleware)
	employeeRouter.Use(googleAuthMiddleware)
	earRouter.Use(googleAuthMiddleware)
	productRouter.Use(googleAuthMiddleware)
	shopRouter.Use(googleAuthMiddleware)
	kafkaRouter.Use(googleAuthMiddleware)

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

func googleAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.String(), "favicon.ico") {
			// Allow favicon.ico to load
			next.ServeHTTP(w, r)
		}

		idToken := r.Header.Get("jwt")
		_, err := utils.ParseJWT(idToken)

		if err != nil {
			utils.WriteErrToClient(w, err)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func getConnectionPool(conf DbConfig) *pgxpool.Pool {
	poolConfig, _ := pgxpool.ParseConfig(conf.DatabaseURL)

	connection, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to database.")

	return connection
}

func getProductController(connpool *pgxpool.Pool, employeeRepo *repositories.EmployeeRepository, kafkaWriter *kafka.Writer, esclient elasticsearch_helpers.ElasticsearchClient) controllers.ProductController {
	productRepository := repositories.NewProductRepo(connpool, kafkaWriter)
	productService := services.ProductService{ProductRepository: productRepository, EmployeeRepository: *employeeRepo}
	productController := controllers.ProductController{Service: productService, ElasticsearchClient: esclient}

	fmt.Println("Product controller up and running.")

	return productController
}

func GetCompanyController(connpool *pgxpool.Pool, kafkaWriter *kafka.Writer) controllers.CompanyController {
	companyRepository := repositories.NewCompanyRepo(connpool, kafkaWriter)
	companyService := services.CompanyService{Repository: companyRepository}
	companyController := controllers.CompanyController{Service: companyService}

	fmt.Println("Company controller up and running.")

	return companyController
}

func getEmployeeController(connpool *pgxpool.Pool) controllers.EmployeeController {
	employeeRepository := repositories.NewEmployeeRepo(connpool)
	employeeService := services.EmployeeService{Repository: employeeRepository}
	employeeController := controllers.EmployeeController{Service: employeeService}

	fmt.Println("\nEmployee controller up and running.")

	return employeeController
}

func getExternalRightController(connpool *pgxpool.Pool) controllers.ExternalRightController {
	earRepository := repositories.NewExternalRightRepo(connpool)
	earService := services.ExternalRightService{Repository: earRepository}
	ExternalRightController := controllers.ExternalRightController{Service: earService}

	fmt.Println("External access rights controller up and running.")

	return ExternalRightController
}

func getConstraintController(connpool *pgxpool.Pool) controllers.ConstraintController {
	constraintRepository := repositories.NewConstraintRepo(connpool)
	constraintService := services.ConstraintService{Repository: constraintRepository}
	constraintController := controllers.ConstraintController{Service: constraintService}

	fmt.Println("Constraints controller up and running.")

	return constraintController
}

func getUserController(connpool *pgxpool.Pool) controllers.UserController {
	userRepository := repositories.NewUserRepo(connpool)
	userService := services.UserService{Repository: userRepository}
	userController := controllers.UserController{Service: userService}

	fmt.Println("User controller up and running.")

	return userController
}

func getShopController(connpool *pgxpool.Pool) controllers.ShopController {
	var conf NominatimConfig
	if _, err := confl.DecodeFile("nominatim_config.conf", &conf); err != nil {
		panic(err)
	}

	shopRepository := repositories.NewShopRepo(connpool)
	geocoder := nominatim.Geocoder(conf.Key)
	shopService := services.ShopService{Repository: shopRepository, Geocoder: geocoder}
	shopController := controllers.ShopController{Service: shopService}

	fmt.Println("Shop controller up and running.")

	return shopController
}

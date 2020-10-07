package repositories

import (
	"context"
	"fmt"
	"internship_project/models"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
)

type config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

var repository ProductRepository

var testProduct models.Product

func TestMain(m *testing.M) {
	connection := instantiateRepository()
	defer connection.Close()

	testProduct = models.Product{
		ID:       "",
		Name:     "TEST_PRODUCT",
		Price:    99,
		Quantity: 10,
		IDC:      "5f46ed8c-03d1-11eb-adc1-0242ac120002",
	}

	ClearTable(connection)
	defer ClearTable(connection)

	os.Exit(m.Run())
}

func instantiateRepository() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("./../database.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.TestDatabaseURL)

	dbtest, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	repository = ProductRepository{DB: dbtest}
	return dbtest
}

func ClearTable(db *pgxpool.Pool) {
	db.Exec(context.Background(),
		"DELETE FROM products")
}

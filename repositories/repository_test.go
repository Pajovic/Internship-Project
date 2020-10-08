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

var repository CompanyRepository

var testCompany models.Company

type config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

func TestMain(m *testing.M) {
	connection := instantiateRepository()
	defer connection.Close()

	testCompany = models.Company{
		Id:     "",
		Name:   "SpaceX",
		IsMain: false,
	}

	clearTable(connection)
	defer clearTable(connection)

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

	repository = CompanyRepository{DB: dbtest}
	return dbtest
}

func clearTable(db *pgxpool.Pool) {
	db.Exec(context.Background(),
		"DELETE FROM companies")
}

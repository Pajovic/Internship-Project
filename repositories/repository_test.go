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

var connection *pgxpool.Pool

type config struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DatabaseURL string `json:"database_url"`
}

func TestMain(m *testing.M) {
	connection = getConnection()
	defer connection.Close()

	ClearTable()
	defer ClearTable()

	os.Exit(m.Run())
}

func getConnection() *pgxpool.Pool {
	var conf config
	if _, err := confl.DecodeFile("./../database_test.conf", &conf); err != nil {
		panic(err)
	}

	poolConfig, _ := pgxpool.ParseConfig(conf.DatabaseURL)

	var err error

	connection, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return connection
}

func ClearTable() {
	connection.Exec(context.Background(),
		"DELETE FROM public.companies")
}

func AddDummyData() {
	company := models.Company{
		Name:   "Facebook",
		IsMain: false,
	}
	AddCompany(&company, connection)
}

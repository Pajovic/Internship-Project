package repositories

import (
	"context"
	"fmt"
	"internship_project/kafka_helpers"
	"internship_project/utils"
	"os"
	"testing"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lytics/confl"
	uuid "github.com/satori/go.uuid"
)

type config struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	DatabaseURL     string `json:"database_url"`
	TestDatabaseURL string `json:"test_database_url"`
}

var (
	Connpool       *pgxpool.Pool
	EmployeeRepo   EmployeeRepository
	ProductRepo    ProductRepository
	CompanyRepo    CompanyRepository
	EarRepo        ExternalRightRepository
	ConstraintRepo ConstraintRepository
)

func TestMain(m *testing.M) {
	Connpool = getConnPool()
	defer Connpool.Close()

	kafkaWriter := kafka_helpers.GetWriter("ava-internship")
	defer kafkaWriter.Close()

	EmployeeRepo = NewEmployeeRepo(Connpool)
	ProductRepo = NewProductRepo(Connpool, kafkaWriter)
	CompanyRepo = NewCompanyRepo(Connpool, kafkaWriter)
	EarRepo = NewExternalRightRepo(Connpool)
	ConstraintRepo = NewConstraintRepo(Connpool)

	utils.SetUpTables(Connpool)

	code := m.Run()

	os.Exit(code)
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

func DoesTableExist(tableName string, connpool *pgxpool.Pool) bool {
	var n int64
	err := connpool.QueryRow(context.Background(), "select 1 from information_schema.tables where table_name=$1", tableName).Scan(&n)
	if err == pgx.ErrNoRows || err != nil {
		return false
	}
	return true
}

func getConnPool() *pgxpool.Pool {
	var conf config
	if _, configerr := confl.DecodeFile("./../dbconfig.conf", &conf); configerr != nil {
		panic(configerr)
	}

	poolConfig, poolerr := pgxpool.ParseConfig(conf.TestDatabaseURL)
	if poolerr != nil {
		panic("Error configuring pool")
	}

	dbtest, dberr := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if dberr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dberr)
		os.Exit(1)
	}

	return dbtest
}

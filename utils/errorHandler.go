package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
)

var (
	NoDataError *pgconn.PgError = &pgconn.PgError{Code: `02000`, Message: `There is no entity with this ID`}
)

// GetErrorMsg is used for error handling
func GetErrorMsg(err error) (string, int) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err.Error(), 400
	}
	switch pgErr.Code {
	case "22P02":
		// invalid ID
		return "You have entered an invalid UUID. Please try again.", 400
	case "23503":
		// foreign key violation
		return "This record can’t be deleted because another record refers to it.", 400
	case "23505":
		// unique constraint violation
		return "This record contains duplicated data that conflicts with what is already in the database.", 400
	case "23514":
		// check constraint violation
		return "This record contains inconsistent or out-of-range data inside column.", 400
	case "22001":
		// value too long for field
		return "This record contains value which exceeds its allowed length.", 400
	case "42P02":
		// invalid parameters
		return "This record contains invalid parametres. " + pgErr.Detail, 400
	case "42601":
		// syntax error
		return "There is a following syntax error in the query:" + "\n" + pgErr.Message, 500
	case "02000":
		// No data
		return pgErr.Message, 404
	case "42P01":
		// Undefined table
		tableName := strings.Split(err.Error(), "\"")[1]
		return "The table you wish to work with, " + tableName + ", does not exist.", 500

	default:
		msg := pgErr.Message
		if d := pgErr.Detail; d != "" {
			msg += "\n\n" + d
		}
		if h := pgErr.Hint; h != "" {
			msg += "\n\n" + h
		}
		return msg, 400
	}
}

// WriteErrToClient is used to return err to client
func WriteErrToClient(w http.ResponseWriter, err error) {
	errMsg, code := GetErrorMsg(err)
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}

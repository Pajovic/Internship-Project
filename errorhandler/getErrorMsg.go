package errorhandler

import (
	"errors"

	"github.com/jackc/pgconn"
)

// GetErrorMsg is used for error handling
func GetErrorMsg(err error) (string, int) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
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
	} else {
		return err.Error(), 500
	}
}

// Generated by go-postgres-codegen 1
package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const CompaniesInsertSql = `
	INSERT INTO 
		public.companies
	(
		id,
		name,
		ismain
	)
	VALUES
		($1,$2,$3)
`

const CompaniesUpdateSql = `
	UPDATE 
		public.companies
	SET
		id=$1,
		name=$2,
		ismain=$3
	WHERE
		id=$4
`

const CompaniesDeleteSql = `
	DELETE FROM
		public.companies
	WHERE
		id=$1
`

type Companies struct {
	Id     pgtype.UUID `db:"id"`
	Name   string      `db:"name"`
	Ismain bool        `db:"ismain"`
}

func (self *Companies) InsertTx(tx *pgx.Tx) (int64, error) {
	commandTag, err := (*tx).Exec(context.Background(), CompaniesInsertSql,
		self.Id,
		self.Name,
		self.Ismain,
	)

	return commandTag.RowsAffected(), err
}

func BatchInsertCompanies(tx *pgx.Tx, batch *[]Companies) (int64, error) {
	vals := []interface{}{}
	stmt := `
	INSERT INTO 
		public.companies
	(
		id,
		name,
		ismain
	)
	VALUES `
	c := 0
	for i, item := range *batch {
		stmt = stmt + fmt.Sprintf(`($%d,$%d,$%d)`, c+1, c+2, c+3)
		if i < len(*batch)-1 {
			stmt = stmt + ","
		}
		vals = append(vals, item.Id, item.Name, item.Ismain)
		c = c + 3
	}

	commandTag, err := (*tx).Exec(context.Background(), stmt, vals...)

	return commandTag.RowsAffected(), err
}

func StrBatchInsertCompanies(batchSize int) string {
	stmt := `
	INSERT INTO 
		public.companies
	(
		id,
		name,
		ismain
	)
	VALUES `
	c := 0
	for i := 0; i < batchSize; i++ {
		stmt = stmt + fmt.Sprintf(`($%d,$%d,$%d)`, c+1, c+2, c+3)
		if i < batchSize-1 {
			stmt = stmt + ","
		}
		c = c + 3
	}
	return stmt
}

func (self *Companies) UpdateTx(tx *pgx.Tx) (int64, error) {
	commandTag, err := (*tx).Exec(context.Background(), CompaniesUpdateSql,
		self.Id,
		self.Name,
		self.Ismain,
		self.Id,
	)

	return commandTag.RowsAffected(), err
}

func (self *Companies) DeleteTx(tx *pgx.Tx) (int64, error) {
	commandTag, err := (*tx).Exec(context.Background(), CompaniesDeleteSql, self.Id)

	return commandTag.RowsAffected(), err
}

func (self *Companies) Scan(rows *pgx.Rows, extensions ...PersistenceExtension) {
	vals, _ := (*rows).Values()
	for i, f := range (*rows).FieldDescriptions() {
		val := vals[i]
		switch string(f.Name) {
		case "id":

			temp := val.([16]uint8)
			uuidVal := pgtype.UUID{}
			uuidVal.Set(temp)
			self.Id = uuidVal

		case "name":
			self.Name = val.(string)
		case "ismain":
			self.Ismain = val.(bool)
		default:
			for _, extension := range extensions {
				extension.Extend(string(f.Name), val)
			}
		}
	}
}
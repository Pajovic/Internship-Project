package repositories

import (
	"context"
	"errors"
	"fmt"
	"internship_project/models"

	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

// GetAllEmployees .
func GetAllEmployees(connection *pgx.Conn) ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	rows, err := connection.Query(context.Background(), "select * from employees")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.CompanyID, &employee.C, &employee.R, &employee.U, &employee.D)
		if err != nil {
			return nil, err
		}
		allEmployees = append(allEmployees, employee)
	}
	return allEmployees, nil
}

// GetEmployeeByID .
func GetEmployeeByID(connection *pgx.Conn, id string) (models.Employee, error) {
	var employee models.Employee
	err := connection.QueryRow(context.Background(), "select * from employees where id=$1", id).Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.CompanyID, &employee.C, &employee.R, &employee.U, &employee.D)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

// AddEmployee .
func AddEmployee(connection *pgx.Conn, employee *models.Employee) error {
	u := uuid.NewV4()
	employee.ID = u.String()
	_, err := connection.Exec(context.Background(), "insert into public.employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		u.Bytes(), employee.FirstName, employee.LastName, employee.CompanyID, employee.C, employee.R, employee.U, employee.D)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpdateEmployee .
func UpdateEmployee(conn *pgx.Conn, updatedEmp models.Employee) error {
	commandTag, err := conn.Exec(context.Background(),
		"UPDATE public.employees SET firstname=$1, lastname=$2, idc=$3, c=$4, r=$5, u=$6, d=$7 WHERE id=$8;",
		updatedEmp.FirstName, updatedEmp.LastName, updatedEmp.CompanyID, updatedEmp.C, updatedEmp.R, updatedEmp.U, updatedEmp.D, updatedEmp.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

// DeleteEmployee .
func DeleteEmployee(conn *pgx.Conn, id string) error {
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM public.employees WHERE id=$1;", id)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

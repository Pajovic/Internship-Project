package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

//EmployeeRepository .
type EmployeeRepository struct {
	DB *pgxpool.Pool
}

// GetAllEmployees .
func (repository *EmployeeRepository) GetAllEmployees(employeeIdc string) ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	query := "select * from employees e where e.idc = $1 or idc in (select * from external_access_rights ear where ear.idrc = $2 and approved = true);"
	rows, err := repository.DB.Query(context.Background(), query, employeeIdc, employeeIdc)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var employee persistence.Employees

		employee.Scan(&rows)

		var employeeUUID string
		err := employee.Id.AssignTo(&employeeUUID)
		if err != nil {
			return allEmployees, err
		}

		var companyUUID string
		err = employee.Idc.AssignTo(&companyUUID)
		if err != nil {
			return allEmployees, err
		}

		allEmployees = append(allEmployees, models.Employee{
			ID:        employeeUUID,
			FirstName: employee.Firstname,
			LastName:  employee.Lastname,
			CompanyID: companyUUID,
			C:         employee.C,
			R:         employee.R,
			U:         employee.U,
			D:         employee.D,
		})
	}
	return allEmployees, nil
}

// GetEmployeeByID .
func (repository *EmployeeRepository) GetEmployeeByID(id string) (models.Employee, error) {
	var employee models.Employee
	rows, err := repository.DB.Query(context.Background(), "select * from employees where id=$1", id)
	defer rows.Close()
	if err != nil {
		return employee, err
	}

	for rows.Next() {
		var employeePers persistence.Employees
		employeePers.Scan(&rows)

		var employeeUUID string
		err := employeePers.Id.AssignTo(&employeeUUID)
		if err != nil {
			return employee, err
		}

		var companyUUID string
		err = employeePers.Idc.AssignTo(&companyUUID)
		if err != nil {
			return employee, err
		}

		employee = models.Employee{
			ID:        employeeUUID,
			FirstName: employeePers.Firstname,
			LastName:  employeePers.Lastname,
			CompanyID: companyUUID,
			C:         employeePers.C,
			R:         employeePers.R,
			U:         employeePers.U,
			D:         employeePers.D,
		}
		break
	}

	return employee, nil
}

// AddEmployee .
func (repository *EmployeeRepository) AddEmployee(employee *models.Employee) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	employeePers := persistence.Employees{
		Firstname: employee.FirstName,
		Lastname:  employee.LastName,
		C:         employee.C,
		R:         employee.R,
		U:         employee.U,
		D:         employee.D,
	}
	employeePers.Idc.Set(employee.CompanyID)
	employeePers.Id.Set(uuid.NewV4())

	_, err = employeePers.InsertTx(&tx)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil
}

// UpdateEmployee .
func (repository *EmployeeRepository) UpdateEmployee(employee models.Employee) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	employeePers := persistence.Employees{
		Firstname: employee.FirstName,
		Lastname:  employee.LastName,
		C:         employee.C,
		R:         employee.R,
		U:         employee.U,
		D:         employee.D,
	}
	employeePers.Idc.Set(employee.CompanyID)
	employeePers.Id.Set(employee.ID)

	commandTag, err := employeePers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to update")
	}

	tx.Commit(context.Background())
	return nil
}

// DeleteEmployee .
func (repository *EmployeeRepository) DeleteEmployee(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	employeePers := persistence.Employees{}
	employeePers.Id.Set(id)

	commandTag, err := employeePers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

// GetEmployeeExternalPermissions .
func (repository *EmployeeRepository) GetEmployeeExternalPermissions(idReceivingCompany string, product models.Product) (models.ExternalRights, error) {
	var allRights []models.ExternalRights
	var rights models.ExternalRights

	// 1. Using idReceivingCompany and idSharingCompany, acquire all external access rules for these two companies
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return rights, err
	}

	defer tx.Rollback(context.Background())

	queryExternalAccess := "SELECT * FROM external_access_rights WHERE idrc = $1 AND idsc = $2;"
	rows, err := tx.Query(context.Background(), queryExternalAccess, idReceivingCompany, product.IDC)

	if err != nil {
		return rights, err
	}
	for rows.Next() {
		var right models.ExternalRights
		err := rows.Scan(&right.ID, &right.IDSC, &right.IDRC, &right.Read, &right.Update, &right.Delete, &right.Approved)
		if err != nil {
			return rights, err
		}
		allRights = append(allRights, right)
	}
	rows.Close()

	if len(allRights) == 0 {
		// 2a. If there is no rows returned employees from this company can't interact with this product
		return rights, errors.New("You don't have any permission for this product")
	} else if len(allRights) == 1 {
		// 2b. If there is only 1 row returned, it means there are no constraints and we can return safely
		rights = allRights[0]
	} else {
		// 2c. Otherwise, we need to acquire constraints, using ID of all constraints
		for _, right := range allRights {
			var accessConstraint models.AccessConstraint
			err := tx.QueryRow(context.Background(), "select * from access_constraints where idear=$1", right.ID).
				Scan(&accessConstraint.ID, &accessConstraint.IDEAR, &accessConstraint.OperatorID, &accessConstraint.PropertyID, &accessConstraint.PropertyValue)
			if err != nil {
				return rights, err
			}
			if checkConstraint(accessConstraint, product) && right.Approved {
				rights = right
			}
		}
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return rights, err
	}

	return rights, nil
}

// CheckCompaniesSharingEmployeeData .
func (repository *EmployeeRepository) CheckCompaniesSharingEmployeeData(idReceivingCompany string, idSharingCompany string) (bool, error) {
	var allRights []models.ExternalRights

	// 1. Using idReceivingCompany and idSharingCompany, acquire all external access rules for these two companies
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return false, err
	}

	defer tx.Rollback(context.Background())

	queryExternalAccess := "SELECT * FROM external_access_rights WHERE idrc = $1 AND idsc = $2;"
	rows, err := tx.Query(context.Background(), queryExternalAccess, idReceivingCompany, idSharingCompany)

	if err != nil {
		return false, err
	}
	for rows.Next() {
		var right models.ExternalRights
		err := rows.Scan(&right.ID, &right.IDSC, &right.IDRC, &right.Read, &right.Update, &right.Delete, &right.Approved)
		if err != nil {
			return false, err
		}
		allRights = append(allRights, right)
	}
	rows.Close()

	if len(allRights) == 0 {
		// 2a. If there is no rows returned employees from this company can't see employees from other companies
		return false, errors.New("Your company does not have rights needed")
	}
	// 2b. If there is any sharing right, we must check if the sharing has been approved
	rightsApproved := false
	for _, right := range allRights {
		if right.Approved {
			rightsApproved = true
			break
		}
	}

	if !rightsApproved {
		return false, errors.New("Sharing between your companies has not been approved")
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return false, err
	}

	return true, nil
}

func checkConstraint(accessConstraint models.AccessConstraint, product models.Product) bool {
	var quantity int32 = int32(accessConstraint.PropertyValue)
	if accessConstraint.OperatorID == 1 {
		if product.Quantity > quantity {
			return true
		} else {
			return false
		}
	} else if accessConstraint.OperatorID == 2 {
		if product.Quantity >= quantity {
			return true
		} else {
			return false
		}
	} else if accessConstraint.OperatorID == 3 {
		if product.Quantity < quantity {
			return true
		} else {
			return false
		}
	} else {
		if product.Quantity <= quantity {
			return true
		} else {
			return false
		}
	}
}

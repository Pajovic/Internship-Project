package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type EmployeeRepository interface {
	GetAllEmployees(string) ([]models.Employee, error)
	GetEmployeeByID(id string) (models.Employee, error)
	AddEmployee(*models.Employee) error
	UpdateEmployee(models.Employee) error
	DeleteEmployee(string) error
	GetEmployeeExternalPermissions(string, models.Product) (models.ExternalRights, error)
	CheckCompaniesSharingEmployeeData(string, string) error
}

type employeeRepository struct {
	DB *pgxpool.Pool
}

func NewEmployeeRepo(db *pgxpool.Pool) EmployeeRepository {
	if db == nil {
		panic("EmployeeRepository not created, pgxpool is nil")
	}
	return &employeeRepository {
		DB: db,
	}
}

// GetAllEmployees .
func (repository *employeeRepository) GetAllEmployees(employeeIdc string) ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	query := "select * from employees e where e.idc = $1 or idc in (select idsc from external_access_rights ear where ear.idrc = $2 and approved = true);"
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
func (repository *employeeRepository) GetEmployeeByID(id string) (models.Employee, error) {
	var employee models.Employee

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return employee, err
	}

	rows, err := repository.DB.Query(context.Background(), "select * from employees where id=$1", Uuid)
	defer rows.Close()

	if err != nil {
		return employee, err
	}

	if !rows.Next() {
		return employee, errors.New("There is no employee with this id")
	}

	var employeePers persistence.Employees
	employeePers.Scan(&rows)

	var employeeUUID string
	err = employeePers.Id.AssignTo(&employeeUUID)
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

	return employee, nil
}

// AddEmployee .
func (repository *employeeRepository) AddEmployee(employee *models.Employee) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	employee.ID = uuid.NewV4().String()

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

	_, err = employeePers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

// UpdateEmployee .
func (repository *employeeRepository) UpdateEmployee(employee models.Employee) error {
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
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

// DeleteEmployee .
func (repository *employeeRepository) DeleteEmployee(id string) error {
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
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

// GetEmployeeExternalPermissions .
func (repository *employeeRepository) GetEmployeeExternalPermissions(idReceivingCompany string, product models.Product) (models.ExternalRights, error) {
	allRights := []models.ExternalRights{}
	var rights models.ExternalRights

	// 1. Using idReceivingCompany and idSharingCompany, acquire all external access rules for these two companies
	queryExternalAccess := "SELECT * FROM external_access_rights WHERE idrc = $1 AND idsc = $2;"
	rows, err := repository.DB.Query(context.Background(), queryExternalAccess, idReceivingCompany, product.IDC)
	defer rows.Close()

	if err != nil {
		return rights, err
	}
	for rows.Next() {
		var ear persistence.ExternalAccessRights
		ear.Scan(&rows)

		var stringUUID string
		err := ear.Id.AssignTo(&stringUUID)
		if err != nil {
			return rights, err
		}

		var idscUUID string
		err = ear.Idsc.AssignTo(&idscUUID)
		if err != nil {
			return rights, err
		}

		var idrcUUID string
		err = ear.Idrc.AssignTo(&idrcUUID)
		if err != nil {
			return rights, err
		}

		allRights = append(allRights, models.ExternalRights{
			ID:       stringUUID,
			IDSC:     idscUUID,
			IDRC:     idrcUUID,
			Read:     ear.R,
			Update:   ear.U,
			Delete:   ear.D,
			Approved: ear.Approved,
		})
	}

	rightsFound := false
	switch len(allRights) {
	case 1:
		// 2a. If there is only 1 row returned, it means there are no constraints and we can return safely
		if allRights[0].Approved {
			rights = allRights[0]
			rightsFound = true
		}
	default:
		// 2b. Otherwise, we need to acquire constraints, using ID of all constraints
		for _, right := range allRights {
			if !right.Approved {
				continue
			}
			rows, err := repository.DB.Query(context.Background(), `select * from access_constraints where idear = $1`, right.ID)
			defer rows.Close()

			if err != nil {
				return rights, err
			}

			rows.Next()
			var constraintPers persistence.AccessConstraints
			constraintPers.Scan(&rows)

			var stringUUID string
			err = constraintPers.Id.AssignTo(&stringUUID)
			if err != nil {
				return rights, err
			}

			var idearUUID string
			err = constraintPers.Idear.AssignTo(&idearUUID)
			if err != nil {
				return rights, err
			}

			var constraint models.AccessConstraint = models.AccessConstraint{
				ID:            stringUUID,
				IDEAR:         idearUUID,
				OperatorID:    constraintPers.OperatorId,
				PropertyID:    constraintPers.PropertyId,
				PropertyValue: constraintPers.PropertyValue,
			}

			if checkConstraint(constraint, product) {
				rights = right
				rightsFound = true
				break
			}
		}
	}

	if !rightsFound {
		return rights, errors.New("Your company does not have rights needed")
	}

	return rights, nil
}

// CheckCompaniesSharingEmployeeData .
func (repository *employeeRepository) CheckCompaniesSharingEmployeeData(idReceivingCompany, idSharingCompany string) error {
	allRights := []models.ExternalRights{}

	// 1. Using idReceivingCompany and idSharingCompany, acquire all external access rules for these two companies

	queryExternalAccess := "SELECT * FROM external_access_rights WHERE idrc = $1 AND idsc = $2 AND approved = true;"
	rows, err := repository.DB.Query(context.Background(), queryExternalAccess, idReceivingCompany, idSharingCompany)
	defer rows.Close()

	if err != nil {
		return err
	}
	for rows.Next() {
		var ear persistence.ExternalAccessRights
		ear.Scan(&rows)

		var stringUUID string
		err := ear.Id.AssignTo(&stringUUID)
		if err != nil {
			return err
		}

		var idscUUID string
		err = ear.Idsc.AssignTo(&idscUUID)
		if err != nil {
			return err
		}

		var idrcUUID string
		err = ear.Idrc.AssignTo(&idrcUUID)
		if err != nil {
			return err
		}

		allRights = append(allRights, models.ExternalRights{
			ID:       stringUUID,
			IDSC:     idscUUID,
			IDRC:     idrcUUID,
			Read:     ear.R,
			Update:   ear.U,
			Delete:   ear.D,
			Approved: ear.Approved,
		})
	}

	if len(allRights) == 0 {
		// 2a. If there is no rows returned employees from this company can't see employees from other companies
		return errors.New("Your company does not have rights needed")
	}

	return nil
}

func checkConstraint(accessConstraint models.AccessConstraint, product models.Product) bool {
	var quantity int32 = int32(accessConstraint.PropertyValue)
	switch accessConstraint.OperatorID {
	case 1:
		return product.Quantity > quantity
	case 2:
		return product.Quantity >= quantity
	case 3:
		return product.Quantity < quantity
	default:
		return product.Quantity <= quantity
	}
}

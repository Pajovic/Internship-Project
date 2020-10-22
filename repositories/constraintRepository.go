package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ConstraintRepository struct {
	DB *pgxpool.Pool
}

func (repository *ConstraintRepository) GetAllConstraints() ([]models.AccessConstraint, error) {
	var constraints []models.AccessConstraint = []models.AccessConstraint{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.access_constraints")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var constraint persistence.AccessConstraints
		constraint.Scan(&rows)

		var stringUUID string
		err = constraint.Id.AssignTo(&stringUUID)
		if err != nil {
			return constraints, err
		}

		var idearUUID string
		err = constraint.Idear.AssignTo(&idearUUID)
		if err != nil {
			return constraints, err
		}

		constraints = append(constraints, models.AccessConstraint{
			ID:            stringUUID,
			IDEAR:         idearUUID,
			OperatorID:    constraint.OperatorId,
			PropertyID:    constraint.PropertyId,
			PropertyValue: constraint.PropertyValue,
		})
	}
	return constraints, nil
}

func (repository *ConstraintRepository) GetConstraint(id string) (models.AccessConstraint, error) {
	var constraint models.AccessConstraint

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return constraint, err
	}

	rows, err := repository.DB.Query(context.Background(), `select * from access_constraints where id = $1`, Uuid)
	defer rows.Close()

	if err != nil {
		return constraint, err
	}

	if !rows.Next() {
		return constraint, errors.New("There is no constraint with this id")
	}

	var constraintPers persistence.AccessConstraints
	constraintPers.Scan(&rows)

	var stringUUID string
	err = constraintPers.Id.AssignTo(&stringUUID)
	if err != nil {
		return constraint, err
	}

	var idearUUID string
	err = constraintPers.Idear.AssignTo(&idearUUID)
	if err != nil {
		return constraint, err
	}

	constraint = models.AccessConstraint{
		ID:            stringUUID,
		IDEAR:         idearUUID,
		OperatorID:    constraintPers.OperatorId,
		PropertyID:    constraintPers.PropertyId,
		PropertyValue: constraintPers.PropertyValue,
	}

	return constraint, nil
}

func (repository *ConstraintRepository) AddConstraint(constraint *models.AccessConstraint) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	constraint.ID = uuid.NewV4().String()
	constraintPers := persistence.AccessConstraints{
		OperatorId:    constraint.OperatorID,
		PropertyId:    constraint.PropertyID,
		PropertyValue: constraint.PropertyValue,
	}
	constraintPers.Id.Set(constraint.ID)
	constraintPers.Idear.Set(constraint.IDEAR)

	_, err = constraintPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *ConstraintRepository) UpdateConstraint(constraint models.AccessConstraint) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	constraintPers := persistence.AccessConstraints{
		OperatorId:    constraint.OperatorID,
		PropertyId:    constraint.PropertyID,
		PropertyValue: constraint.PropertyValue,
	}
	constraintPers.Id.Set(constraint.ID)
	constraintPers.Idear.Set(constraint.IDEAR)

	commandTag, err := constraintPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to update")
	}

	return tx.Commit(context.Background())
}

func (repository *ConstraintRepository) DeleteConstraint(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	constraintPers := persistence.AccessConstraints{}
	constraintPers.Id.Set(id)

	commandTag, err := constraintPers.DeleteTx(&tx)

	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to delete")
	}
	return tx.Commit(context.Background())
}

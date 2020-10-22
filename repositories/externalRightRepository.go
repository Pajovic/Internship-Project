package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ExternalRightRepository struct {
	DB *pgxpool.Pool
}

func (repository *ExternalRightRepository) GetAllEars() ([]models.ExternalRights, error) {
	var ears []models.ExternalRights = []models.ExternalRights{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.external_access_rights")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ear persistence.ExternalAccessRights
		ear.Scan(&rows)

		var stringUUID string
		err := ear.Id.AssignTo(&stringUUID)
		if err != nil {
			return ears, err
		}

		var idscUUID string
		err = ear.Idsc.AssignTo(&idscUUID)
		if err != nil {
			return ears, err
		}

		var idrcUUID string
		err = ear.Idrc.AssignTo(&idrcUUID)
		if err != nil {
			return ears, err
		}

		ears = append(ears, models.ExternalRights{
			ID:       stringUUID,
			IDSC:     idscUUID,
			IDRC:     idrcUUID,
			Read:     ear.R,
			Update:   ear.U,
			Delete:   ear.D,
			Approved: ear.Approved,
		})
	}
	return ears, nil
}

func (repository *ExternalRightRepository) GetEar(id string) (models.ExternalRights, error) {
	var ear models.ExternalRights

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return ear, err
	}

	rows, err := repository.DB.Query(context.Background(), `select * from external_access_rights where id = $1`, Uuid)
	defer rows.Close()

	if err != nil {
		return ear, err
	}

	if !rows.Next() {
		return ear, errors.New("There is no ear with this id")
	}

	var earPers persistence.ExternalAccessRights
	earPers.Scan(&rows)

	var stringUUID string
	err = earPers.Id.AssignTo(&stringUUID)
	if err != nil {
		return ear, err
	}

	var idscUUID string
	err = earPers.Idsc.AssignTo(&idscUUID)
	if err != nil {
		return ear, err
	}

	var idrcUUID string
	err = earPers.Idrc.AssignTo(&idrcUUID)
	if err != nil {
		return ear, err
	}

	ear = models.ExternalRights{
		ID:       stringUUID,
		IDSC:     idscUUID,
		IDRC:     idrcUUID,
		Read:     earPers.R,
		Update:   earPers.U,
		Delete:   earPers.D,
		Approved: earPers.Approved,
	}

	return ear, nil
}

func (repository *ExternalRightRepository) AddEar(ear *models.ExternalRights) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	ear.ID = uuid.NewV4().String()
	earPers := persistence.ExternalAccessRights{
		R:        ear.Read,
		U:        ear.Update,
		D:        ear.Delete,
		Approved: ear.Approved,
	}
	earPers.Id.Set(ear.ID)
	earPers.Idsc.Set(ear.IDSC)
	earPers.Idrc.Set(ear.IDRC)

	_, err = earPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *ExternalRightRepository) UpdateEar(ear models.ExternalRights) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	earPers := persistence.ExternalAccessRights{
		R:        ear.Read,
		U:        ear.Update,
		D:        ear.Delete,
		Approved: ear.Approved,
	}
	earPers.Id.Set(ear.ID)
	earPers.Idsc.Set(ear.IDSC)
	earPers.Idrc.Set(ear.IDRC)

	commandTag, err := earPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to update")
	}

	return tx.Commit(context.Background())
}

func (repository *ExternalRightRepository) DeleteEar(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	earPers := persistence.ExternalAccessRights{}
	earPers.Id.Set(id)

	commandTag, err := earPers.DeleteTx(&tx)

	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to delete")
	}
	return tx.Commit(context.Background())
}

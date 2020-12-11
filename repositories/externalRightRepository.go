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

type ExternalRightRepository interface {
	GetAllEars() ([]models.ExternalRights, error)
	GetEar(id string) (models.ExternalRights, error)
	AddEar(ear *models.ExternalRights) error
	UpdateEar(ear models.ExternalRights) error
	DeleteEar(id string) error
	DeleteExternalRightsForCompany(string) error
}

type externalRightRepository struct {
	DB              *pgxpool.Pool
	ConstraintsRepo ConstraintRepository
}

func NewExternalRightRepo(db *pgxpool.Pool) ExternalRightRepository {
	if db == nil {
		panic("ExternalRightRepository not created, pgxpool is nil")
	}
	return &externalRightRepository{
		DB:              db,
		ConstraintsRepo: NewConstraintRepo(db),
	}
}

func (repository *externalRightRepository) GetAllEars() ([]models.ExternalRights, error) {
	ears := []models.ExternalRights{}
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

func (repository *externalRightRepository) GetEar(id string) (models.ExternalRights, error) {
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

func (repository *externalRightRepository) AddEar(ear *models.ExternalRights) error {
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

func (repository *externalRightRepository) UpdateEar(ear models.ExternalRights) error {
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
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *externalRightRepository) DeleteEar(id string) error {
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
		return utils.NoDataError
	}
	return tx.Commit(context.Background())
}

func (repository *externalRightRepository) DeleteExternalRightsForCompany(idc string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `DELETE FROM external_access_rights WHERE idsc=$1 or idrc = $1`

	_, err = tx.Exec(context.Background(), query, idc)

	if err != nil {
		return err
	}

	err = repository.ConstraintsRepo.DeleteConstraintsForCompany(idc)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}

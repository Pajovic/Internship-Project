package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"
)

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUser(string) (models.User, error)
	AddUser(models.User) error
	UpdateUser(models.User) error
	DeleteUser(string) error
}

type userRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	if db == nil {
		panic("UserRepository not created, pgxpool is nil")
	}
	return &userRepository {
		DB: db,
	}
}

func (repository *userRepository) GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.users")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user persistence.Users
		user.Scan(&rows)

		users = append(users, models.User{
			ID:     user.Id,
			Email: user.Email,
			Name:   user.Name,
		})
	}
	return users, nil
}

func (repository *userRepository) GetUser(id string) (models.User, error) {
	var user models.User

	rows, err := repository.DB.Query(context.Background(), `select * from public.users where id = $1`, id)
	defer rows.Close()

	if err != nil {
		return user, err
	}

	if !rows.Next() {
		return user, utils.NoDataError
	}

	var userPers persistence.Users
	userPers.Scan(&rows)

	user = models.User{
		ID:     userPers.Id,
		Email:  userPers.Email,
		Name:   userPers.Name,
	}

	return user, nil
}

func (repository *userRepository) AddUser(user models.User) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	userPers := persistence.Users{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	_, err = userPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *userRepository) UpdateUser(user models.User) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	userPers := persistence.Users{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	commandTag, err := userPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *userRepository) DeleteUser(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	userPers := persistence.Users{}
	userPers.Id = id

	commandTag, err := userPers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}
	return tx.Commit(context.Background())
}



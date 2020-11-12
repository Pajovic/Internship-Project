package services

import (
	"errors"
	"internship_project/models"
	"internship_project/repositories"

	"github.com/markbates/goth"
)

type UserService struct {
	Repository repositories.UserRepository
}

func (service *UserService) GetUser(id string) (models.User, error) {
	user, err := service.Repository.GetUser(id)

	if err != nil {
		return models.User{}, errors.New("You are not authenticated. Please sign in with Google to continue")
	}

	return user, nil
}

func (service *UserService) GoogleSignIn(u goth.User) (models.User, error) {
	var user models.User
	exists, err := service.Repository.DoesUserExists(u.UserID)
	if err != nil {
		return user, err
	}
	if !exists {
		user = models.User{
			ID:    u.UserID,
			Email: u.Email,
			Name:  u.FirstName + " " + u.LastName,
		}
		err = service.Repository.AddUser(user)
	} else {
		user, err = service.Repository.GetUser(u.UserID)
	}
	if err != nil {
		return user, err
	}
	return user, nil
}

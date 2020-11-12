package services

import (
	"github.com/markbates/goth"
	"internship_project/models"
	"internship_project/repositories"
)

type UserService struct {
	Repository repositories.UserRepository
}

func (service *UserService) GoogleSignIn(u goth.User) (models.User, error) {
	var user models.User
	exists, err := service.Repository.DoesUserExists(u.UserID)
	if err != nil {
		return user, err
	}
	if !exists {
		user = models.User{
			ID:     u.UserID,
			Email: 	u.Email,
			Name:   u.FirstName + " " + u.LastName,
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

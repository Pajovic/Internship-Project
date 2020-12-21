package services

import (
	"internship_project/models"
	"internship_project/repositories"
	"internship_project/utils"
)

type UserService struct {
	Repository repositories.UserRepository
}

func (service *UserService) GoogleSignIn(token string) (models.User, error) {
	var user models.User

	claims, err := utils.ValidateGoogleJWT(token)
	if err != nil {
		return user, err
	}

	exists, err := service.Repository.DoesUserExists(claims.Sub)
	if err != nil {
		return user, err
	}
	if !exists {
		user = models.User{
			ID:    claims.Sub,
			Email: claims.Email,
			Name:  claims.FirstName + " " + claims.LastName,
		}
		err = service.Repository.AddUser(user)
	} else {
		user, err = service.Repository.GetUser(claims.Sub)
	}
	if err != nil {
		return user, err
	}
	return user, nil
}


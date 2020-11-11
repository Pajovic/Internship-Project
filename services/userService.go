package services

import (
	"context"
	"internship_project/models"
	"internship_project/repositories"

	"google.golang.org/api/idtoken"
)

type UserService struct {
	Repository repositories.UserRepository
}

const googleClientId = ""

func (service *UserService) GoogleSignIn(token string) (models.User, error) {
	var user models.User

	tokenValidator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		return user, err
	}

	payload, err := tokenValidator.Validate(context.Background(), token, googleClientId)
	if err != nil {
		return user, err
	}

	sub := payload.Subject
	user, err = service.Repository.GetUser(sub)
	if err != nil {
		user = models.User{
			ID:     sub,
			Email: "",
			Name:   "",
		}
		service.Repository.AddUser(user)
	}

	return user, nil
}

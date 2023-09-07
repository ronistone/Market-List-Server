package service

import (
	"context"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUser(ctx context.Context, id int64) (model.User, error)
}
type User struct {
	UserRepository repository.UserRepository
}

func CreateUserService(userRepository repository.UserRepository) UserService {
	return &User{
		UserRepository: userRepository,
	}
}

func (u User) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	return u.UserRepository.CreateUser(ctx, user)
}

func (u User) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := u.UserRepository.GetUserById(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

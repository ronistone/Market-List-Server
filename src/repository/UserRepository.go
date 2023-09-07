package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/ronistone/market-list/src/model"
	"github.com/ronistone/market-list/src/util"
)

type UserRepository interface {
	CreateUser(ctx context.Context, purchase model.User) (model.User, error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
}

type User struct {
	DbConnection *dbr.Connection
}

func CreateUserRepository(connection *dbr.Connection) UserRepository {
	return &User{
		DbConnection: connection,
	}
}

func (p User) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	INSERT INTO market_user(name, password) 
		VALUES (?, ?) 
	RETURNING *
	`, user.Name, user.Password)

	_, err := statement.LoadContext(ctx, &user)
	if err != nil {
		return model.User{}, util.MakeErrorUnknown(err)
	}

	return user, nil
}

func (p User) GetUserById(ctx context.Context, id int64) (model.User, error) {
	statement := p.DbConnection.NewSession(nil).SelectBySql(`
	SELECT * FROM MARKET_USER where ID = ?
	`, id)

	var user model.User
	err := statement.LoadOne(&user)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return model.User{}, util.MakeError(util.NOT_FOUND, fmt.Sprintf("User %d not found", id))
		}
		return model.User{}, util.MakeErrorUnknown(err)
	}

	return user, nil
}

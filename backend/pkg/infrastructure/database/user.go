package database

import (
	"github.com/ansel1/merry"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/google/uuid"
	"github.com/joseenriqe97/test-cabri/config"
	"github.com/joseenriqe97/test-cabri/pkg/infrastructure/domain"
)

type userDatabase struct {
	db            *goqu.Database
	userTableGoqu exp.IdentifierExpression
}

type UserRepositoryInteface interface {
	Create(user *domain.User) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
}

var UserDatabase UserRepositoryInteface = &userDatabase{
	db:            &config.SQLDBGoqu,
	userTableGoqu: goqu.T("user"),
}

func (u *userDatabase) Create(user *domain.User) (*domain.User, error) {
	user.ID = uuid.New().String()

	_, err := u.db.Insert(u.userTableGoqu).Cols(
		"id",
		"name",
		"last_name",
		"email",
		"password",
	).Vals(goqu.Vals{
		user.ID,
		user.Name,
		user.LastName,
		user.Email,
		user.Password,
	}).Executor().Exec()
	if err != nil {
		return nil, merry.Wrap(err)
	}
	return user, nil
}

func (u *userDatabase) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	_, err := u.db.From(u.userTableGoqu).Where(
		u.userTableGoqu.Col("email").Eq(email),
	).ScanStruct(&user)

	if err != nil {
		return nil, merry.Wrap(err)
	}

	return &user, nil
}
func (u *userDatabase) GetByID(id string) (*domain.User, error) {
	var user domain.User

	_, err := u.db.From(u.userTableGoqu).Where(
		u.userTableGoqu.Col("id").Eq(id),
	).ScanStruct(&user)

	if err != nil {
		return nil, merry.Wrap(err)
	}

	// if !found {
	// 	return nil, merry.New(fmt.Sprintf("user %s not found", id))
	// }

	return &user, nil
}

package repository

import (
	"github.com/azka-zaydan/article-materials/unit-testing/user/model"
	"github.com/jmoiron/sqlx"
)

//go:generate go run go.uber.org/mock/mockgen -source=./repo.go -destination=../mocks/repository_mock.go -package=mocks

type UserRepository interface {
	FindUserByID(id int) (res model.User, err error)
	FindUserByEmail(email string) (res model.User, err error)
	CreateUser(user *model.User) (err error)
	DoesUserExist(email string) (exist bool, err error)
}

type UserRepositoryImpl struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) FindUserByID(id int) (res model.User, err error) {
	err = r.DB.Get(&res, "SELECT * FROM users WHERE id = ?", id)
	return
}

func (r *UserRepositoryImpl) FindUserByEmail(email string) (res model.User, err error) {
	err = r.DB.Get(&res, "SELECT * FROM users WHERE email = ?", email)
	return
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) (err error) {
	_, err = r.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	return
}

func (r *UserRepositoryImpl) DoesUserExist(email string) (exist bool, err error) {
	var count int
	err = r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ?", email)
	if err != nil {
		return
	}
	return count > 0, nil
}

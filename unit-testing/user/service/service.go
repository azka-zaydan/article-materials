package service

import (
	"database/sql"
	"errors"

	"github.com/azka-zaydan/article-materials/unit-testing/user/model"
	"github.com/azka-zaydan/article-materials/unit-testing/user/model/dto"
	"github.com/azka-zaydan/article-materials/unit-testing/user/repository"
)

//go:generate go run go.uber.org/mock/mockgen -source=./service.go -destination=../mocks/service_mock.go -package=mocks

type UserService interface {
	GetUserByID(id int) (res model.User, err error)
	GetUserByEmail(email string) (res model.User, err error)
	CreateUser(req dto.CreateUserReq) (err error)
}

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

func (s *UserServiceImpl) GetUserByID(id int) (res model.User, err error) {
	res, err = s.UserRepo.FindUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.New("user not found")
		}
		err = errors.New("internal server error")
		return
	}
	return
}

func (s *UserServiceImpl) GetUserByEmail(email string) (res model.User, err error) {
	res, err = s.UserRepo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.New("user not found")
		}
		err = errors.New("internal server error")
		return
	}
	return
}

func (s *UserServiceImpl) CreateUser(req dto.CreateUserReq) (err error) {
	user := model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	exist, err := s.UserRepo.DoesUserExist(user.Email)
	if err != nil {
		return
	}
	if exist {
		return errors.New("user already exist")
	}

	err = s.UserRepo.CreateUser(&user)
	return
}

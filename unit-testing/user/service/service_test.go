package service_test

import (
	"database/sql"
	"testing"

	"github.com/azka-zaydan/article-materials/unit-testing/user/mocks"
	"github.com/azka-zaydan/article-materials/unit-testing/user/model"
	"github.com/azka-zaydan/article-materials/unit-testing/user/model/dto"
	"github.com/azka-zaydan/article-materials/unit-testing/user/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserServiceImpl_GetUserByID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := service.NewUserService(mockUserRepo)

	userMock := model.User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByID(1).Return(userMock, nil)
		res, err := service.GetUserByID(1)

		assert.NoError(t, err)
		assert.Equal(t, userMock, res)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByID(1).Return(model.User{}, assert.AnError)
		res, err := service.GetUserByID(1)

		assert.Error(t, err)
		assert.Equal(t, model.User{}, res)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByID(1).Return(model.User{}, sql.ErrNoRows)
		res, err := service.GetUserByID(1)

		assert.Error(t, err)
		assert.Equal(t, model.User{}, res)
	})
}

func TestUserServiceImpl_GetUserByEmail(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := service.NewUserService(mockUserRepo)

	userMock := model.User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}
	johnEmail := "john@example.com"

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByEmail(johnEmail).Return(userMock, nil)
		res, err := service.GetUserByEmail(johnEmail)

		assert.NoError(t, err)
		assert.Equal(t, userMock, res)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByEmail(johnEmail).Return(model.User{}, assert.AnError)
		res, err := service.GetUserByEmail(johnEmail)

		assert.Error(t, err)
		assert.Equal(t, model.User{}, res)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().FindUserByEmail(johnEmail).Return(model.User{}, sql.ErrNoRows)
		res, err := service.GetUserByEmail(johnEmail)

		assert.Error(t, err)
		assert.Equal(t, model.User{}, res)
	})
}

func TestUserServiceImpl_CreateUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := service.NewUserService(mockUserRepo)

	createUserReq := dto.CreateUserReq{
		Name:  "John",
		Email: "john@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().DoesUserExist(createUserReq.Email).Return(false, nil)
		mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
		err := service.CreateUser(createUserReq)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.EXPECT().DoesUserExist(createUserReq.Email).Return(false, assert.AnError)
		mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(assert.AnError)
		err := service.CreateUser(createUserReq)

		assert.Error(t, err)
	})

	t.Run("user already exist", func(t *testing.T) {
		mockUserRepo.EXPECT().DoesUserExist(createUserReq.Email).Return(true, nil)
		err := service.CreateUser(createUserReq)

		assert.Error(t, err)
	})

}

package service

import (
	"context"
	"errors"

	"go-web-demo/internal/common/response"
	"go-web-demo/internal/dal"
	"go-web-demo/internal/db"
	"go-web-demo/internal/dto"
	"go-web-demo/internal/model"

	"gorm.io/gorm"
)

type UserService struct{}

var userService *UserService

func init() {
	userService = &UserService{}
}

func GetUserService() *UserService {
	return userService
}

func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest) (*model.User, error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Status:   req.Status,
	}

	if err := dal.GetUserDao().Insert(db.DB, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *dto.UserUpdateRequest) (*model.User, error) {
	user, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	user.Gender = req.Gender
	user.Status = req.Status

	if err := dal.GetUserDao().UpdateSelective(db.DB, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return dal.GetUserDao().DeleteSoft(db.DB, int64(user.ID))
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := dal.GetUserDao().SelectById(db.DB, int64(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetList(ctx context.Context, req *dto.UserPageQueryRequest) ([]model.User, error) {
	return dal.GetUserDao().SelectList(db.DB, req)
}

func (s *UserService) GetPageList(ctx context.Context, req *dto.UserPageQueryRequest) (*response.PageResult[model.User], error) {
	return dal.GetUserDao().SelectListByPage(db.DB, req)
}
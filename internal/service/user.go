package service 

import (
	"context"
	"errors"
	"strings"

	"service-xp/internal/model"
	"service-xp/internal/repo"

)

var ErrInvalidEmail = errors.New("invalid email")

var ErrInvalidName = errors.New("invalid name")

type UserService struct{
	repo *repo.UserRepo
}


func NewUserService(r *repo.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, name string, email string) (*model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetUser(ctx context.Context, email string) (*model.User, error){
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, ErrInvalidEmail
	}
	return s.repo.GetByEmail(ctx, email)
}


// 2.service — метод CreateUser: обрезать пробелы,
//  привести email к нижнему регистру, проверить что name и email не пустые. 
//  Вернуть ошибку если пусто.
package repo

import (
	"context"
	"database/sql"
	"errors"

	"service-xp/internal/model"
)

// 1.repo — метод Create, который делает INSERT и получает id обратно. 
// SQL — только здесь.

var ErrNotFound = errors.New("not found")

type UserRepo struct{
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}


func (r *UserRepo) Create(ctx context.Context, name, email string) (*model.User, error) {
	const q = `insert into public.users (name, email) values ($1, $2) returning id`
	u := &model.User{Name: name, Email: email}// create user User u = new User(name, email);  
	err := r.db.QueryRowContext(ctx, q, name, email).Scan(&u.ID)
	return u, err
}


func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User,error){
	const q = `SELECT id, name, email FROM get_user_by_email($1)`

	var u model.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
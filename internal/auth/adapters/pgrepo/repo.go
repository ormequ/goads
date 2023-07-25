package pgrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/auth/app"
	"goads/internal/auth/users"
)

type Repo struct {
	db *pgx.Conn
}

const insertUserQuery = `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`

func (r Repo) Store(ctx context.Context, user users.User) (int64, error) {
	var id int64 = -1
	err := r.db.QueryRow(ctx, insertUserQuery, user.Email, user.Name, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_email_key" {
				err = errors.Join(err, app.ErrEmailAlreadyExists)
			}
		}
	}
	return id, err
}

const getUserByEmailQuery = `SELECT id, email, password FROM users WHERE email=$1`

func (r Repo) GetByEmail(ctx context.Context, email string) (users.User, error) {
	var usr users.User
	err := r.db.QueryRow(ctx, getUserByEmailQuery, email).Scan(&usr.ID, &usr.Email, &usr.Password)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrIncorrectCredentials)
	}
	return usr, err
}

const getUserByIdQuery = `SELECT id, email, name FROM users WHERE id=$1`

func (r Repo) GetByID(ctx context.Context, id int64) (users.User, error) {
	if ctx.Err() != nil {
		return users.User{}, ctx.Err()
	}
	var user users.User
	err := r.db.QueryRow(ctx, getUserByIdQuery, id).Scan(&user.ID, &user.Email, &user.Name)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return user, err
}

const updateUserQuery = `UPDATE users SET email=$1, name=$2, password=$3 WHERE id=$4`

func (r Repo) Update(ctx context.Context, user users.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(ctx, updateUserQuery, user.Email, user.Name, user.Password, user.ID)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

const deleteUserQuery = `DELETE FROM users WHERE id=$1`

func (r Repo) Delete(ctx context.Context, id int64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(ctx, deleteUserQuery, id)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func New(conn *pgx.Conn) Repo {
	return Repo{db: conn}
}

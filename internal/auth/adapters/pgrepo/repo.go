package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/auth/app"
	"goads/internal/auth/users"
	"goads/internal/pkg/errwrap"
)

type Repo struct {
	db *pgx.Conn
}

const (
	constrEmail = "users_email_key"
)

func (r Repo) Store(ctx context.Context, user users.User) (int64, error) {
	const query = `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`
	const op = "pgrepo.Store"

	var id int64 = -1
	err := r.db.QueryRow(ctx, query, user.Email, user.Name, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrEmail {
			err = errwrap.New(app.ErrEmailAlreadyExists, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
	}
	return id, err
}

func (r Repo) GetByEmail(ctx context.Context, email string) (users.User, error) {
	const query = `SELECT id, email, password FROM users WHERE email=$1`
	const op = "pgrepo.GetByEmail"

	var usr users.User
	err := r.db.QueryRow(ctx, query, email).Scan(&usr.ID, &usr.Email, &usr.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrIncorrectCredentials, app.ServiceName, op).
			WithDetails(fmt.Sprintf("%v | email: %s", err.Error(), email)).
			OnObject("user", usr.ID)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).
			WithDetails(fmt.Sprintf("email: %s", email)).
			OnObject("user", usr.ID)
	}
	return usr, err
}

func (r Repo) GetByID(ctx context.Context, id int64) (users.User, error) {
	const query = `SELECT id, email, name FROM users WHERE id=$1`
	const op = "pgrepo.GetByID"

	var user users.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("user", id)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("user", id)
	}
	return user, err
}

func (r Repo) Update(ctx context.Context, user users.User) error {
	const query = `UPDATE users SET email=$1, name=$2, password=$3 WHERE id=$4`
	const op = "pgrepo.Update"

	_, err := r.db.Exec(ctx, query, user.Email, user.Name, user.Password, user.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("user", user.ID)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("user", user.ID)
	}
	return err
}

func (r Repo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM users WHERE id=$1`
	const op = "pgrepo.Delete"

	_, err := r.db.Exec(ctx, query, id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("user", id)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("user", id)
	}
	return err
}

func New(conn *pgx.Conn) Repo {
	return Repo{db: conn}
}

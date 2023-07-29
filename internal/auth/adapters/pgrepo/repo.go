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

const (
	constrEmail = "users_email_key"
)

func (r Repo) Store(ctx context.Context, user users.User) (int64, error) {
	const query = `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`

	var id int64 = -1
	err := r.db.QueryRow(ctx, query, user.Email, user.Name, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == constrEmail {
				err = errors.Join(err, app.ErrEmailAlreadyExists)
			}
		}
	}
	return id, err
}

func (r Repo) GetByEmail(ctx context.Context, email string) (users.User, error) {
	const query = `SELECT id, email, password FROM users WHERE email=$1`

	var usr users.User
	err := r.db.QueryRow(ctx, query, email).Scan(&usr.ID, &usr.Email, &usr.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(app.ErrIncorrectCredentials)
	}
	return usr, err
}

func (r Repo) GetByID(ctx context.Context, id int64) (users.User, error) {
	const query = `SELECT id, email, name FROM users WHERE id=$1`

	var user users.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(app.ErrNotFound)
	}
	return user, err
}

func (r Repo) Update(ctx context.Context, user users.User) error {
	const query = `UPDATE users SET email=$1, name=$2, password=$3 WHERE id=$4`

	_, err := r.db.Exec(ctx, query, user.Email, user.Name, user.Password, user.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func (r Repo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM users WHERE id=$1`

	_, err := r.db.Exec(ctx, query, id)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func New(conn *pgx.Conn) Repo {
	return Repo{db: conn}
}

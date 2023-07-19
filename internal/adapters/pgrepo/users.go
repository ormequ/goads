package pgrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/app"
	"goads/internal/entities/users"
)

type Users struct {
	db *pgx.Conn
}

func (r Users) Store(ctx context.Context, user users.User) (int64, error) {
	var id int64 = -1
	if ctx.Err() != nil {
		return id, ctx.Err()
	}
	err := r.db.QueryRow(ctx, `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`, user.Email, user.Name).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_email_key" {
				err = errors.Join(err, app.ErrAlreadyExists)
			}
		}
	}
	return id, err
}

func (r Users) GetByID(ctx context.Context, id int64) (users.User, error) {
	if ctx.Err() != nil {
		return users.User{}, ctx.Err()
	}
	var user users.User
	err := r.db.QueryRow(ctx, `SELECT id, email, name FROM users WHERE id=$1`, id).Scan(&user.ID, &user.Email, &user.Name)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return user, err
}

func (r Users) GetByEmail(ctx context.Context, email string) (users.User, error) {
	if ctx.Err() != nil {
		return users.User{}, ctx.Err()
	}
	var user users.User
	err := r.db.QueryRow(ctx, `SELECT id, email, name FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Email, &user.Name)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return user, err
}

func (r Users) Update(ctx context.Context, user users.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(ctx, `UPDATE users SET email=$1, name=$2 WHERE id=$3`, user.Email, user.Name, user.ID)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func (r Users) Delete(ctx context.Context, id int64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func NewUsers(conn *pgx.Conn) Users {
	return Users{db: conn}
}

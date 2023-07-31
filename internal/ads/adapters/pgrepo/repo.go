package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/ads/ads"
	"goads/internal/ads/app"
	"goads/internal/pkg/errwrap"
	"time"
)

type Repo struct {
	db *pgx.Conn
}

const (
	constrAuthorID = "ads_author_id_key"
)

func (r Repo) Store(ctx context.Context, ad ads.Ad) (int64, error) {
	const query = `INSERT INTO ads (author_id, published, title, text, create_date, update_date) 	
			       VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	const op = "pgrepo.Store"

	var id int64 = -1
	err := r.db.QueryRow(
		ctx, query,
		ad.AuthorID, ad.Published, ad.Title,
		ad.Text, ad.GetCreateDate(), ad.GetUpdateDate(),
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrAuthorID {
			err = errwrap.New(app.ErrAuthorNotFound, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
	}
	return id, err
}

func (r Repo) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	const query = `SELECT id, author_id, published, title, text, create_date, update_date 
                   FROM ads WHERE id=$1`
	const op = "pgrepo.GetByID"

	var ad ads.Ad
	err := r.db.QueryRow(ctx, query, id).
		Scan(&ad.ID, &ad.AuthorID, &ad.Published, &ad.Title, &ad.Text, &ad.CreateDate, &ad.UpdateDate)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrAdNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("ad", id)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("ad", id)
	}
	return ad, err
}

func (r Repo) GetOnlyPublished(ctx context.Context, ids []int64) (res []ads.Ad, err error) {
	const query = `SELECT id, author_id, published, title, text, create_date, update_date 
			  FROM ads WHERE published = true AND id = ANY($1)`
	const op = "pgrepo.GetOnlyPublished"

	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("given ads: %v", ids))
		}
	}()

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errwrap.New(app.ErrAdNotFound, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
		return
	}

	for rows.Next() {
		var ad ads.Ad
		err = rows.Scan(&ad.ID, &ad.AuthorID, &ad.Published, &ad.Title, &ad.Text, &ad.CreateDate, &ad.UpdateDate)
		if err != nil {
			err = errwrap.New(err, app.ServiceName, op)
			return
		}
		res = append(res, ad)
	}
	return
}

func (r Repo) GetFiltered(ctx context.Context, filter app.Filter) ([]ads.Ad, error) {
	const op = "pgrepo.GetFiltered"

	// query forms dynamically
	query := "SELECT id, author_id, published, title, text, create_date, update_date FROM ads"

	where := false
	if !filter.All {
		query += " WHERE published=true"
		where = true
	}
	if filter.AuthorID != -1 {
		if !where {
			query += " WHERE "
			where = true
		} else {
			query += " AND "
		}
		// %d cannot create sql injection
		query += fmt.Sprintf("author_id=%d", filter.AuthorID)
	}
	if !filter.Date.IsZero() && filter.Date.UnixMilli() != 0 {
		if !where {
			query += " WHERE "
			where = true
		} else {
			query += " AND "
		}
		// %d cannot create sql injection
		d := filter.Date.Truncate(time.Hour * 24)
		query += fmt.Sprintf("create_date='%d-%d-%d'", d.Year(), d.Month(), d.Day())
	}
	if !where {
		query += " WHERE "
	} else {
		query += " AND "
	}
	query += "title LIKE '%' || $1 || '%'"
	rows, err := r.db.Query(ctx, query, filter.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errwrap.New(app.ErrAdNotFound, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
		return nil, err

	}
	defer rows.Close()

	var res []ads.Ad
	for rows.Next() {
		ad := ads.Ad{}
		err := rows.Scan(&ad.ID, &ad.AuthorID, &ad.Published, &ad.Title, &ad.Text, &ad.CreateDate, &ad.UpdateDate)
		if err != nil {
			return res, errwrap.New(err, app.ServiceName, op)
		}
		res = append(res, ad)
	}
	return res, nil
}

func (r Repo) Update(ctx context.Context, ad ads.Ad) error {
	const query = `UPDATE ads SET author_id=$1, published=$2, title=$3, text=$4, create_date=$5, update_date=$6 
                   WHERE id=$7`
	const op = "pgrepo.Update"

	_, err := r.db.Exec(
		ctx, query,
		ad.AuthorID, ad.Published, ad.Title, ad.Text, ad.GetCreateDate(), ad.GetUpdateDate(), ad.ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrAdNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("ad", ad.ID)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("ad", ad.ID)
	}
	return err
}

func (r Repo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM ads WHERE id=$1`
	const op = "pgrepo.Delete"

	_, err := r.db.Exec(ctx, query, id)

	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrAdNotFound, app.ServiceName, op).WithDetails(err.Error()).OnObject("ad", id)
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("ad", id)
	}
	return err
}

func New(conn *pgx.Conn) Repo {
	return Repo{db: conn}
}

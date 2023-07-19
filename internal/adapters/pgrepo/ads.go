package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/app"
	ads2 "goads/internal/app/ad"
	"goads/internal/entities/ads"
	"time"
)

type Ads struct {
	db *pgx.Conn
}

func (r Ads) Store(ctx context.Context, ad ads.Ad) (int64, error) {
	var id int64 = -1
	if ctx.Err() != nil {
		return id, ctx.Err()
	}
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO ads (author_id, published, title, text, create_date, update_date) 	
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		ad.AuthorID, ad.Published, ad.Title, ad.Text, ad.GetCreateDate(), ad.GetUpdateDate(),
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "ads_author_id_key" {
				err = errors.Join(err, app.ErrNotFound)
			}
		}
	}
	return id, err
}

func (r Ads) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	if ctx.Err() != nil {
		return ads.Ad{}, ctx.Err()
	}
	var ad ads.Ad
	err := r.db.QueryRow(ctx, `SELECT id, author_id, published, title, text, create_date, update_date FROM ads WHERE id=$1`, id).
		Scan(&ad.ID, &ad.AuthorID, &ad.Published, &ad.Title, &ad.Text, &ad.CreateDate, &ad.UpdateDate)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return ad, err
}

func (r Ads) GetFiltered(ctx context.Context, filter ads2.Filter) ([]ads.Ad, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

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
	query += "title LIKE $1 || '%'"
	rows, err := r.db.Query(ctx, query, filter.Prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ads.Ad
	for rows.Next() {
		ad := ads.Ad{}
		err := rows.Scan(&ad.ID, &ad.AuthorID, &ad.Published, &ad.Title, &ad.Text, &ad.CreateDate, &ad.UpdateDate)
		if err != nil {
			return nil, err
		}
		result = append(result, ad)
	}
	return result, nil
}

func (r Ads) Update(ctx context.Context, ad ads.Ad) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(
		ctx, `UPDATE ads SET author_id=$1, published=$2, title=$3, text=$4, create_date=$5, update_date=$6 WHERE id=$7`,
		ad.AuthorID, ad.Published, ad.Title, ad.Text, ad.GetCreateDate(), ad.GetUpdateDate(), ad.ID,
	)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func (r Ads) Delete(ctx context.Context, id int64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := r.db.Exec(ctx, `DELETE FROM ads WHERE id=$1`, id)
	if err == pgx.ErrNoRows {
		err = errors.Join(app.ErrNotFound)
	}
	return err
}

func NewAds(conn *pgx.Conn) Ads {
	return Ads{db: conn}
}

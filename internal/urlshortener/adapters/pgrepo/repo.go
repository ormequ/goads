package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"goads/internal/pkg/errwrap"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/entities/links"
	"sort"
)

type Repo struct {
	db *pgx.Conn
}

const (
	constrAlias = "links_alias_key"
	constrAdID  = "link_ads_ad_id_key"
)

func (r Repo) SizeApprox(ctx context.Context) (int64, error) {
	const query = `SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = 'db'`
	const op = "pgrepo.SizeApprox"

	var sz int64 = -1
	err := r.db.QueryRow(ctx, query).Scan(&sz)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		err = errwrap.New(err, app.ServiceName, op)
	}
	return sz, err
}

func (r Repo) Store(ctx context.Context, link links.Link) (id int64, err error) {
	const linksQuery = `INSERT INTO links (alias, url, author_id) VALUES ($1, $2, $3) RETURNING id`
	const adsQuery = `INSERT INTO link_ads (link_id, ad_id) VALUES ($1, $2)`
	const op = "pgrepo.Store"

	id = -1
	err = r.db.QueryRow(ctx, linksQuery, link.Alias, link.URL, link.AuthorID).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrAlias {
			err = errwrap.New(app.ErrAlreadyExists, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
		return
	}

	if len(link.Ads) == 0 {
		return
	}

	batchAds := &pgx.Batch{}
	for _, adID := range link.Ads {
		batchAds.Queue(adsQuery, id, adID)
	}
	br := r.db.SendBatch(ctx, batchAds)
	defer func() {
		if err != nil {
			err = errors.Join(err, br.Close())
		} else if err = br.Close(); err != nil {
			err = errwrap.New(err, app.ServiceName, op)
		}
	}()
	_, err = br.Exec()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrAdID {
			err = errwrap.New(app.ErrAdNotExists, app.ServiceName, op).WithDetails(err.Error())
		} else {
			err = errwrap.New(err, app.ServiceName, op)
		}
	}
	return
}

func pgArrayToGoSlice(arr pgtype.Array[pgtype.Int8]) []int64 {
	if !arr.Valid {
		return nil
	}
	res := make([]int64, 0, len(arr.Elements))
	for i := range arr.Elements {
		if arr.Elements[i].Valid {
			res = append(res, arr.Elements[i].Int64)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}

func (r Repo) GetByID(ctx context.Context, id int64) (links.Link, error) {
	const query = `
		SELECT links.id, alias, url, author_id, ARRAY_AGG(la.ad_id ORDER BY la.ad_id)
		FROM links
        	LEFT JOIN link_ads la on links.id = la.link_id 
		WHERE links.id=$1
		GROUP BY links.id;
	`
	const op = "pgrepo.GetByID"

	link := links.Link{}
	var ads pgtype.Array[pgtype.Int8]

	err := r.db.QueryRow(ctx, query, id).Scan(&link.ID, &link.Alias, &link.URL, &link.AuthorID, &ads)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errwrap.New(app.ErrNotFound, app.ServiceName, op).
				WithDetails(err.Error()).
				OnObject("link", id)
		} else {
			err = errwrap.New(err, app.ServiceName, op).
				OnObject("link", id)
		}
		return link, err
	}
	link.Ads = pgArrayToGoSlice(ads)
	return link, nil
}

func (r Repo) GetByAuthor(ctx context.Context, authorID int64) ([]links.Link, error) {
	const query = `
		SELECT id, alias, url, ARRAY_AGG(la.ad_id ORDER BY la.ad_id)
		FROM links 
			LEFT JOIN link_ads la on links.id=la.link_id 
		WHERE links.author_id=$1
		GROUP BY links.id
	`
	const op = "pgrepo.GetByAuthor"

	rows, err := r.db.Query(ctx, query, authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errwrap.New(app.ErrNotFound, app.ServiceName, op).
				WithDetails(fmt.Sprintf("%v | author ID: %d", err, authorID))
		} else {
			err = errwrap.New(err, app.ServiceName, op).
				WithDetails(fmt.Sprintf("author ID: %d", authorID))
		}
		return nil, err
	}
	defer rows.Close()
	var res []links.Link
	for rows.Next() {
		link := links.Link{AuthorID: authorID}
		var ads pgtype.Array[pgtype.Int8]
		err := rows.Scan(&link.ID, &link.Alias, &link.URL, &ads)
		if err != nil {
			return res, errwrap.New(err, app.ServiceName, op)
		}
		link.Ads = pgArrayToGoSlice(ads)
		res = append(res, link)
	}
	return res, nil
}

func (r Repo) GetByAlias(ctx context.Context, alias string) (links.Link, error) {
	const query = `SELECT links.id, url, author_id, ARRAY_AGG(la.ad_id ORDER BY la.ad_id)
		FROM links
				 LEFT JOIN link_ads la on links.id = la.link_id
		WHERE links.alias = $1
		GROUP BY links.id
	`
	const op = "pgrepo.GetByAlias"

	link := links.Link{Alias: alias}
	var ads pgtype.Array[pgtype.Int8]
	err := r.db.QueryRow(ctx, query, alias).Scan(&link.ID, &link.URL, &link.AuthorID, &ads)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errwrap.New(app.ErrNotFound, app.ServiceName, op).
				WithDetails(fmt.Sprintf("%v | alias: %s", err, alias))
		} else {
			err = errwrap.New(err, app.ServiceName, op).
				WithDetails(fmt.Sprintf("alias: %s", alias))
		}
		return link, err
	}
	link.Ads = pgArrayToGoSlice(ads)
	return link, nil
}

func (r Repo) UpdateAlias(ctx context.Context, id int64, alias string) error {
	const query = `UPDATE links SET alias=$1 WHERE id=$2`
	const op = "pgrepo.UpdateAlias"

	_, err := r.db.Exec(ctx, query, alias, id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).OnObject("link", id).WithDetails(err.Error())
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("link", id)
	}
	return err
}

func (r Repo) AddAd(ctx context.Context, linkID int64, adID int64) error {
	const query = `INSERT INTO link_ads (link_id, ad_id) VALUES ($1, $2)`
	const op = "pgrepo.AddAd"

	_, err := r.db.Exec(ctx, query, linkID, adID)
	if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("link", linkID).
			WithDetails(fmt.Sprintf("ad ID: %d", adID))
	}
	return err
}

func (r Repo) DeleteAd(ctx context.Context, linkID int64, adID int64) error {
	const query = `DELETE FROM link_ads WHERE link_id=$1 AND ad_id=$2`
	const op = "pgrepo.DeleteAd"

	_, err := r.db.Exec(ctx, query, linkID, adID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).OnObject("link", linkID).
			WithDetails(fmt.Sprintf("%v | ad ID: %d", err, adID))
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).OnObject("link", linkID).
			WithDetails(fmt.Sprintf("ad ID: %d", adID))
	}
	return err
}

func (r Repo) Delete(ctx context.Context, id int64) error {
	const linksQuery = `DELETE FROM links WHERE id=$1`
	const op = "pgrepo.Delete"

	_, err := r.db.Exec(ctx, linksQuery, id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errwrap.New(app.ErrNotFound, app.ServiceName, op).
			OnObject("link", id).
			WithDetails(err.Error())
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op).
			OnObject("link", id)
	}
	return err
}

func New(db *pgx.Conn) Repo {
	return Repo{db}
}

package pgrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/links"
)

type Repo struct {
	links *pgx.Conn
	ads   *pgx.Conn
}

const (
	constrAlias = "links_alias_key"
	constrAdID  = "link_ads_ad_id_key"
)

func checkConstraint(err error, constr string) bool {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.ConstraintName == constr
		}
	}
	return false
}

func (r Repo) SizeApprox(ctx context.Context) (int64, error) {
	const query = `SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = 'links'`
	var sz int64 = -1
	err := r.links.QueryRow(ctx, query).Scan(&sz)
	return sz, err
}

func (r Repo) Store(ctx context.Context, link links.Link) (int64, error) {
	const linksQuery = `INSERT INTO links (alias, url, author_id) VALUES ($1, $2, $3) RETURNING id`
	const adsQuery = `INSERT INTO link_ads (link_id, ad_id) VALUES ($1, $2)`

	var id int64 = -1
	err := r.links.QueryRow(ctx, linksQuery, link.Alias, link.URL, link.AuthorID).Scan(&id)
	if checkConstraint(err, constrAlias) {
		return -1, errors.Join(err, app.ErrAlreadyExists)
	}
	if err != nil {
		return -1, err
	}

	if len(link.Ads) == 0 {
		return id, nil
	}

	batchAds := &pgx.Batch{}
	for _, adID := range link.Ads {
		batchAds.Queue(adsQuery, id, adID)
	}
	br := r.ads.SendBatch(ctx, batchAds)
	defer func() { err = errors.Join(err, br.Close()) }()
	_, err = br.Exec()
	if checkConstraint(err, constrAdID) {
		err = errors.Join(err, app.ErrAdNotExists)
	}

	return id, err
}

func (r Repo) getAdsByLink(ctx context.Context, linkID int64) ([]int64, error) {
	const query = `SELECT ad_id FROM link_ads WHERE link_id=$1`

	rows, err := r.ads.Query(ctx, query, linkID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Join(app.ErrNoAds)
	}
	var res []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return res, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r Repo) GetByID(ctx context.Context, id int64) (links.Link, error) {
	const query = `SELECT alias, url, author_id FROM links WHERE id=$1`

	link := links.Link{ID: id}
	err := r.links.QueryRow(ctx, query, id).Scan(&link.Alias, &link.URL, &link.AuthorID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	if err != nil {
		return link, err
	}

	link.Ads, err = r.getAdsByLink(ctx, link.ID)
	return link, err
}

func (r Repo) getAdsByLinks(ctx context.Context, ids []int64) (map[int64][]int64, error) {
	const query = `SELECT link_id, ad_id FROM link_ads WHERE link_id = ANY($1)`

	rows, err := r.ads.Query(ctx, query, ids)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Join(app.ErrNoAds)
	}
	res := make(map[int64][]int64)
	for rows.Next() {
		var link, ad int64
		err := rows.Scan(&link, &ad)
		if err != nil {
			return res, err
		}
		res[link] = append(res[link], ad)
	}
	return res, nil
}

func (r Repo) GetByAuthor(ctx context.Context, authorID int64) ([]links.Link, error) {
	const query = `SELECT id, alias, url FROM links WHERE author_id=$1`

	rows, err := r.links.Query(ctx, query, authorID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Join(app.ErrNotFound)
	}
	lMap := make(map[int64]links.Link)
	var lIDs []int64
	for rows.Next() {
		link := links.Link{AuthorID: authorID}
		err := rows.Scan(&link.ID, &link.Alias, &link.URL)
		if err != nil {
			return nil, err
		}
		lMap[link.ID] = link
		lIDs = append(lIDs, link.ID)
	}

	ads, err := r.getAdsByLinks(ctx, lIDs)
	if err != nil {
		return nil, err
	}

	var res []links.Link
	for id, link := range lMap {
		if ads != nil {
			link.Ads = ads[id]
		}
		res = append(res, link)
	}
	return res, nil
}

func (r Repo) GetByAlias(ctx context.Context, alias string) (links.Link, error) {
	const query = `SELECT id, url, author_id FROM links WHERE alias=$1`

	link := links.Link{Alias: alias}
	err := r.links.QueryRow(ctx, query, alias).Scan(&link.ID, &link.URL, &link.AuthorID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	if err != nil {
		return link, err
	}

	link.Ads, err = r.getAdsByLink(ctx, link.ID)
	return link, err
}

func (r Repo) UpdateAlias(ctx context.Context, id int64, alias string) error {
	const query = `UPDATE links SET alias=$1 WHERE id=$2`

	_, err := r.links.Exec(ctx, query, alias, id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	return err
}

func (r Repo) AddAd(ctx context.Context, linkID int64, adID int64) error {
	const query = `INSERT INTO link_ads (link_id, ad_id) VALUES ($1, $2)`

	_, err := r.links.Exec(ctx, query, linkID, adID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	return err
}

func (r Repo) DeleteAd(ctx context.Context, linkID int64, adID int64) error {
	const query = `DELETE FROM link_ads WHERE link_id=$1 AND ad_id=$2`

	_, err := r.links.Exec(ctx, query, linkID, adID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	return err
}

func (r Repo) Delete(ctx context.Context, id int64) error {
	const linksQuery = `DELETE FROM links WHERE id=$1`
	const adsQuery = `DELETE FROM link_ads WHERE link_id=$1`

	_, err := r.links.Exec(ctx, linksQuery, id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.Join(err, app.ErrNotFound)
	}
	if err == nil {
		_, err = r.ads.Exec(ctx, adsQuery, id)
		if errors.Is(err, pgx.ErrNoRows) {
			err = errors.Join(err, app.ErrNoAds)
		}
	}
	return err
}

func New(links *pgx.Conn, ads *pgx.Conn) Repo {
	return Repo{links, ads}
}

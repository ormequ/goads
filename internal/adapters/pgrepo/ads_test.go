package pgrepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"testing"
)

func TestAds_GetFiltered(t *testing.T) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:root@localhost:5432/goads")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()
	//urepo := NewUsers(conn)
	//id, err := urepo.Store(ctx, users.New("no-reply@test.com", "test"))
	//fmt.Println(id)
	//require.NoError(t, err)
	//u, err := urepo.GetByEmail(ctx, "test@test.com")
	//require.NoError(t, err)
	//repo := NewAds(conn)
	////require.NoError(t, repo.Store(ctx, ads.New("title`';\"", "\"';`", u.ID)))
	////require.NoError(t, repo.Store(ctx, ads.New("title of interesting ad", "interesting", u.ID)))
	////require.NoError(t, repo.Store(ctx, ads.New("title of ad", "ad", u.ID)))
	//a, err := repo.GetFiltered(ctx, ad.Filter{
	//	AuthorID: u.ID,
	//	Date:     time.Now(),
	//	Prefix:   "title ",
	//	All:      true,
	//})
	//fmt.Println(a)
	//require.NoError(t, err)
}

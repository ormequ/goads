package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"goads/internal/adapters/maprepo"
	"goads/internal/app/providers"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	grpcPort "goads/internal/ports/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"
)

func getTestGRPCClient(t *testing.T) (grpcPort.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		_ = lis.Close()
	})

	srv := grpc.NewServer(grpcPort.GetUnaryInterceptors())
	t.Cleanup(func() {
		srv.Stop()
	})

	u := providers.NewUsers(maprepo.New[users.User]())
	a := providers.NewAds(maprepo.New[ads.Ad](), u)
	svc := grpcPort.NewService(a, u)
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		_ = conn.Close()
	})
	return grpcPort.NewAdServiceClient(conn), ctx
}

func TestGRRPCCreateUser(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@test.com"})
	assert.NoError(t, err)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "test@test.com", res.Email)
	assert.Equal(t, int64(0), res.Id)
}

func TestGRPCUpdateUser(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	user, err = client.ChangeUserEmail(ctx, &grpcPort.ChangeUserEmailRequest{Id: user.Id, Email: "asdf@asdf.com"})
	assert.NoError(t, err)
	assert.Equal(t, "asdf@asdf.com", user.Email)

	user, err = client.ChangeUserName(ctx, &grpcPort.ChangeUserNameRequest{Id: user.Id, Name: "asdf"})
	assert.NoError(t, err)
	assert.Equal(t, "asdf", user.Name)
}

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)
	assert.Zero(t, ad.Id)
	assert.Equal(t, "hello", ad.Title)
	assert.Equal(t, "world", ad.Text)
	assert.Equal(t, user.Id, ad.AuthorId)
	assert.False(t, ad.Published)
	assert.WithinDuration(t, time.Now().UTC(), time.UnixMilli(ad.CreateDate), time.Millisecond*2)
	assert.WithinDuration(t, time.Now().UTC(), time.UnixMilli(ad.CreateDate), time.Millisecond*2)
}

func TestGRPCCreateAdWithoutUser(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 123, Title: "hello", Text: "world"})
	code := status.Code(err)
	assert.Equal(t, codes.NotFound, code)
}

func TestGRPCIncorrectGet(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	_, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 123})
	code := status.Code(err)
	assert.Equal(t, codes.NotFound, code)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)
	time.Sleep(5 * time.Millisecond)
	ad, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, ad.Published)
	assert.WithinDuration(t, time.Now().UTC(), time.UnixMilli(ad.UpdateDate), time.Millisecond*2)
	assert.NotEqual(t, time.UnixMilli(ad.CreateDate).Truncate(time.Millisecond), time.UnixMilli(ad.UpdateDate).Truncate(time.Millisecond))
}

func TestGRPCUpdateAd(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ad, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: user.Id, AdId: ad.Id, Title: "привет", Text: "мир"})
	assert.NoError(t, err)
	assert.Equal(t, "привет", ad.Title)
	assert.Equal(t, "мир", ad.Text)
}

func TestGRPCListAds(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ad, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user.Id, AdId: ad.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "best cat", Text: "not for sale"})
	assert.NoError(t, err)

	adsResp, err := client.ListAds(ctx, &grpcPort.FilterAdsRequest{AuthorId: -1})
	listAds := adsResp.List
	assert.NoError(t, err)
	assert.Len(t, listAds, 1)
	assert.Equal(t, ad.Id, listAds[0].Id)
	assert.Equal(t, ad.Title, listAds[0].Title)
	assert.Equal(t, ad.Text, listAds[0].Text)
	assert.Equal(t, ad.AuthorId, listAds[0].AuthorId)
	assert.True(t, listAds[0].Published)
}

func TestGRPCSearch(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: " hello", Text: "world"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello world", Text: "world"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello world ", Text: "world"})
	assert.NoError(t, err)

	l1, err := client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Title: "hello"})
	assert.NoError(t, err)
	assert.Len(t, l1.List, 3)

	l2, err := client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Title: " hello"})
	assert.NoError(t, err)
	assert.Len(t, l2.List, 1)
	assert.Equal(t, ad.Id, l2.List[0].Id)

	l3, err := client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Title: "hello world"})
	assert.NoError(t, err)
	assert.Len(t, l3.List, 2)
}

func TestGRPCDeleteAd(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	adsList, err := client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Title: "hello"})
	assert.NoError(t, err)
	assert.Len(t, adsList.List, 1)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AuthorId: user.Id, AdId: ad.Id})
	assert.NoError(t, err)

	adsList, err = client.SearchAds(ctx, &grpcPort.SearchAdsRequest{Title: "hello"})
	assert.NoError(t, err)
	assert.Len(t, adsList.List, 0)
}

func TestGRPCDeleteUser(t *testing.T) {
	client, ctx := getTestGRPCClient(t)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "test", Email: "test@test.com"})
	assert.NoError(t, err)

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: user.Id})
	assert.NoError(t, err)

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: user.Id})
	code := status.Code(err)
	assert.Equal(t, codes.NotFound, code)
}

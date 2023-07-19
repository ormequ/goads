package tests

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	require.NoError(t, err)
	assert.Zero(t, user.Data.ID)
	assert.Equal(t, "test@test.com", user.Data.Email)
	assert.Equal(t, "test", user.Data.Name)

	gotUser, err := client.getUser(user.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, gotUser, user)
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	user, err = client.changeUser(user.Data.ID, "email", "asdf@asdf.com")
	assert.NoError(t, err)
	assert.Equal(t, "asdf@asdf.com", user.Data.Email)

	user, err = client.changeUser(user.Data.ID, "name", "asdf")
	assert.NoError(t, err)
	assert.Equal(t, "asdf", user.Data.Name)
}

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	ad, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, ad.Data.ID)
	assert.Equal(t, "hello", ad.Data.Title)
	assert.Equal(t, "world", ad.Data.Text)
	assert.Equal(t, user.Data.ID, ad.Data.AuthorID)
	assert.False(t, ad.Data.Published)
	assert.WithinDuration(t, time.Now().UTC(), ad.Data.CreateDate, time.Millisecond*2)
	assert.WithinDuration(t, time.Now().UTC(), ad.Data.UpdateDate, time.Millisecond*2)
}

func TestCreateAdWithoutUser(t *testing.T) {
	client := getTestClient()
	_, err := client.createAd(123, "hello", "world")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestIncorrectGet(t *testing.T) {
	client := getTestClient()

	_, err := client.getUser(123)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	ad, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.CreateDate.Truncate(time.Millisecond), ad.Data.UpdateDate.Truncate(time.Millisecond))

	time.Sleep(time.Millisecond * 5)
	ad, err = client.changeAdStatus(user.Data.ID, ad.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, ad.Data.Published)
	assert.WithinDuration(t, time.Now().UTC(), ad.Data.UpdateDate, time.Millisecond*2)

	time.Sleep(time.Millisecond * 5)
	ad, err = client.changeAdStatus(user.Data.ID, ad.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, ad.Data.Published)
	assert.WithinDuration(t, time.Now().UTC(), ad.Data.UpdateDate, time.Millisecond*2)

	time.Sleep(time.Millisecond * 5)
	ad, err = client.changeAdStatus(user.Data.ID, ad.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, ad.Data.Published)
	assert.WithinDuration(t, time.Now().UTC(), ad.Data.UpdateDate, time.Millisecond*2)
	assert.NotEqual(t, ad.Data.CreateDate.Truncate(time.Millisecond), ad.Data.UpdateDate.Truncate(time.Millisecond))
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(user.Data.ID, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, "привет", response.Data.Title)
	assert.Equal(t, "мир", response.Data.Text)
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(user.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, publishedAd.Data.ID, ads.Data[0].ID)
	assert.Equal(t, publishedAd.Data.Title, ads.Data[0].Title)
	assert.Equal(t, publishedAd.Data.Text, ads.Data[0].Text)
	assert.Equal(t, publishedAd.Data.AuthorID, ads.Data[0].AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestSearch(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	require.NoError(t, err)

	ad, err := client.createAd(user.Data.ID, " hello", "world")
	require.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello", "world")
	require.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello world", "world")
	require.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello world ", "world")
	require.NoError(t, err)

	l1, err := client.searchAds("hello")
	require.NoError(t, err)
	assert.Len(t, l1.Data, 3)

	// + is URL decoded space
	l2, err := client.searchAds("+hello")
	require.NoError(t, err)
	assert.Len(t, l2.Data, 1)
	assert.Equal(t, ad.Data, l2.Data[0])

	l3, err := client.searchAds("hello+world")
	require.NoError(t, err)
	assert.Len(t, l3.Data, 2)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	ad, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	ads, err := client.searchAds("hello")

	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)

	err = client.deleteAd(user.Data.ID, ad.Data.ID)
	assert.NoError(t, err)

	ads, err = client.searchAds("hello")

	assert.NoError(t, err)
	assert.Len(t, ads.Data, 0)
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	err = client.deleteUser(user.Data.ID)
	assert.NoError(t, err)

	_, err = client.getUser(user.Data.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

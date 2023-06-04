package tests

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestFilterEmpty(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(user.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.getFilteredAds("")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, publishedAd.Data, ads.Data[0])
}

func TestFilterAll(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.getFilteredAds("all")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestFilterDate(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("test@test.com", "test")
	assert.NoError(t, err)

	ad1, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	ad2, err := client.createAd(user.Data.ID, "go", "lang")
	assert.NoError(t, err)

	ad1, err = client.changeAdStatus(user.Data.ID, ad1.Data.ID, true)
	assert.NoError(t, err)
	ad2, err = client.changeAdStatus(user.Data.ID, ad2.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.getFilteredAds("date=" + time.Now().UTC().Format(time.RFC3339Nano))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ad2.Data, ads.Data[0])

}

func TestFilterAuthor(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)
	u2, err := client.createUser("t2@test.com", "t2")
	assert.NoError(t, err)

	empty, err := client.getFilteredAds("author=" + strconv.Itoa(int(u1.Data.ID)))
	assert.NoError(t, err)
	assert.Len(t, empty.Data, 0)

	ad1, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	ad2, err := client.createAd(u2.Data.ID, "go", "lang")
	assert.NoError(t, err)

	ad1, err = client.changeAdStatus(u1.Data.ID, ad1.Data.ID, true)
	assert.NoError(t, err)
	ad2, err = client.changeAdStatus(u2.Data.ID, ad2.Data.ID, true)
	assert.NoError(t, err)

	res1, err := client.getFilteredAds("author=" + strconv.Itoa(int(u1.Data.ID)))
	assert.NoError(t, err)
	assert.Len(t, res1.Data, 1)
	assert.Equal(t, ad1.Data, res1.Data[0])
	res2, err := client.getFilteredAds("author=" + strconv.Itoa(int(u2.Data.ID)))
	assert.NoError(t, err)
	assert.Len(t, res2.Data, 1)
	assert.Equal(t, ad2.Data, res2.Data[0])
}

func TestManyFilters(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)
	u2, err := client.createUser("t2@test.com", "t2")
	assert.NoError(t, err)

	ad11, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	ad21, err := client.createAd(u2.Data.ID, "go", "lang")
	assert.NoError(t, err)

	time.Sleep(time.Second * 3)

	ad12, err := client.createAd(u1.Data.ID, "HELLO", "WORLD")
	assert.NoError(t, err)
	ad22, err := client.createAd(u2.Data.ID, "GO", "LANG")
	assert.NoError(t, err)

	ad21, err = client.changeAdStatus(u2.Data.ID, ad21.Data.ID, true)
	assert.NoError(t, err)
	ad22, err = client.changeAdStatus(u2.Data.ID, ad22.Data.ID, true)
	assert.NoError(t, err)

	res1, err := client.getFilteredAds("all&date=" + ad22.Data.CreateDate.Format(time.RFC3339Nano))
	assert.NoError(t, err)
	assert.Len(t, res1.Data, 2)

	res2, err := client.getFilteredAds("all&author=" + strconv.Itoa(int(u1.Data.ID)) + "&date=" + ad22.Data.CreateDate.Format(time.RFC3339Nano))
	assert.NoError(t, err)
	assert.Len(t, res2.Data, 1)
	assert.Equal(t, ad12.Data, res2.Data[0])

	res3, err := client.getFilteredAds("author=" + strconv.Itoa(int(u2.Data.ID)) + "&date=" + ad11.Data.CreateDate.Format(time.RFC3339Nano))
	assert.NoError(t, err)
	assert.Len(t, res3.Data, 1)
	assert.Equal(t, ad21.Data, res3.Data[0])
}

func TestIncorrectFilter(t *testing.T) {
	client := getTestClient()
	_, err := client.getFilteredAds("author=123")
	assert.ErrorIs(t, err, ErrBadRequest)
	_, err = client.getFilteredAds("author")
	assert.ErrorIs(t, err, ErrBadRequest)
	_, err = client.getFilteredAds("date=123")
	assert.ErrorIs(t, err, ErrBadRequest)
}

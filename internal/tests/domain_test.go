package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)
	u2, err := client.createUser("t2@test.com", "t2")
	assert.NoError(t, err)

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(u2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}
func TestDeleteAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)
	u2, err := client.createUser("t2@test.com", "t2")
	assert.NoError(t, err)

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	err = client.deleteAd(u2.Data.ID, resp.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)
	u2, err := client.createUser("t2@test.com", "t2")
	assert.NoError(t, err)

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(u2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	u1, err := client.createUser("t1@test.com", "t1")
	assert.NoError(t, err)

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.Data.ID)

	resp, err = client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Data.ID)

	resp, err = client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.Data.ID)
}

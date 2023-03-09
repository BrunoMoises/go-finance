package db

import (
	"context"
	"testing"

	"github.com/BrunoMoises/go-finance/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arq := CreateUserParams{
		Username: util.RandomString(6),
		Password: util.RandomString(12),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arq)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arq.Username, user.Username)
	require.Equal(t, arq.Password, user.Password)
	require.Equal(t, arq.Email, user.Email)
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.NotEmpty(t, user2.CreatedAt)
}

func TestGetUserById(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.NotEmpty(t, user2.CreatedAt)
}

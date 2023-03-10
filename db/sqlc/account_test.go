package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/BrunoMoises/go-finance-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	category1 := createRandomCategory(t)
	arg := CreateAccountParams{
		UserID:      category1.UserID,
		CategoryID:  category1.ID,
		Title:       util.RandomString(12),
		Type:        category1.Type,
		Description: util.RandomString(20),
		Value:       10,
		Date:        time.Now(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, arg.CategoryID, account.CategoryID)
	require.Equal(t, arg.Title, account.Title)
	require.Equal(t, arg.Type, account.Type)
	require.Equal(t, arg.Description, account.Description)
	require.Equal(t, arg.Value, account.Value)

	require.NotEmpty(t, account.Date)
	require.NotEmpty(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.UserID, account2.UserID)
	require.Equal(t, account1.CategoryID, account2.CategoryID)
	require.Equal(t, account1.Title, account2.Title)
	require.Equal(t, account1.Type, account2.Type)
	require.Equal(t, account1.Description, account2.Description)
	require.Equal(t, account1.Value, account2.Value)

	require.NotEmpty(t, account2.Date)
	require.NotEmpty(t, account2.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:          account1.ID,
		Title:       util.RandomString(12),
		Description: util.RandomString(20),
		Value:       15,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Title, account2.Title)
	require.Equal(t, arg.Description, account2.Description)
	require.Equal(t, arg.Value, account2.Value)
	require.NotEmpty(t, account2.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountsParams{
		UserID:      lastAccount.UserID,
		Type:        lastAccount.Type,
		Title:       lastAccount.Title,
		Description: lastAccount.Description,
		Date: sql.NullTime{
			Time:  lastAccount.Date,
			Valid: !lastAccount.Date.IsZero(),
		},
		CategoryID: sql.NullInt32{
			Int32: lastAccount.CategoryID,
			Valid: lastAccount.CategoryID != 0,
		},
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.Equal(t, lastAccount.ID, account.ID)
		require.Equal(t, lastAccount.UserID, account.UserID)
		require.Equal(t, lastAccount.Title, account.Title)
		require.Equal(t, lastAccount.Description, account.Description)
		require.Equal(t, lastAccount.Value, account.Value)
		require.NotEmpty(t, lastAccount.CreatedAt)
		require.NotEmpty(t, lastAccount.Date)
	}
}

func TestListGetReports(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountReportsParams{
		UserID: lastAccount.UserID,
		Type:   lastAccount.Type,
	}

	sumValue, err := testQueries.GetAccountReports(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sumValue)
}

func TestListGetGraph(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountGraphParams{
		UserID: lastAccount.UserID,
		Type:   lastAccount.Type,
	}

	graphValue, err := testQueries.GetAccountGraph(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, graphValue)
}

package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

const (
	dataSourceName = "postgresql://yaroslav:AnySecretPassword!!@postgres1:5432/yaaws?sslmode=disable&timezone=UTC"
)

func TestPing(t *testing.T) {
	connection, err := sql.Open("postgres", dataSourceName)
	require.NoError(t, err)
	defer connection.Close()

	require.NoError(t, connection.Ping())
}
package test

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type PostgreMock struct {
	mock.Mock
}

func (*PostgreMock) Close() error {
	return nil
}

func (m *PostgreMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	argsList := m.Called(query, args)

	return argsList.Get(0).(*sql.Rows), argsList.Error(1)
}

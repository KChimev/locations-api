package models

import "database/sql"

type Postgre interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

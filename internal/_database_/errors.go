package database

import "errors"

var ErrDataUpdateNotAllow = errDataUpdateNotAllow()
var ErrDatabaseNotFound = errDatabaseNotFound()

func errDataUpdateNotAllow() error {
	return errors.New("DataUpdateNotAllowed")
}

func errDatabaseNotFound() error {
	return errors.New("DataUpdateNotAllowed")
}

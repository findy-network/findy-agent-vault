package db

type BaseObject struct {
	ID        string `faker:"uuid_hyphenated"`
	CreatedMs int64  `faker:"created"`
}

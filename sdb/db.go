package sdb

type Db interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}

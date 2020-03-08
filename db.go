package main

import bolt "go.etcd.io/bbolt"

type DB struct {
	*bolt.DB
}

func NewDB(dbName string) (*DB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

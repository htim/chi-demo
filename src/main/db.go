package main

import (
	"github.com/jmoiron/sqlx"
)

func GetLinkById(db *sqlx.DB, id int64) (*Link, error) {
	link := []Link{}
	err := db.Select(&link, "select * from links where id=$1", id)
	if err != nil {
		return nil, err
	} else {
		return &link[0], nil
	}
}
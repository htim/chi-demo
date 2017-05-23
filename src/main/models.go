package main

import "strconv"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

func (u *User) String() string {
	return u.Name + " " + u.Login + " " + u.Login
}

type Link struct {
	ID int64 `db:"id"`
	Link string `db:"link"`
	IsPublic bool `db:"is_public"`
}

func (l *Link) String() string {
	return strconv.FormatInt(l.ID, 10) + " " + l.Link + " " + strconv.FormatBool(l.IsPublic)
}
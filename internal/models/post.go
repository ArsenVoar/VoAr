package models

type Post struct {
	ID       int
	UserID   int
	Title    string
	Anons    string
	FullText string
}

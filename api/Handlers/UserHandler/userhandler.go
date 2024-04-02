package userhandler

import "database/sql"

type UserHandler struct {
	conn *sql.DB
}

const (
	LAYOUT = "2006-01-02T15:04:05Z"
)

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{conn: db}
}

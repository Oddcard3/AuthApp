package models

import (
	"authapp/db/models/chats"
	"authapp/db/models/users"
	"database/sql"
)

// SetConn sets DB connection
func SetConn(c *sql.DB) {
	users.SetConn(c)
	chats.SetConn(c)
}

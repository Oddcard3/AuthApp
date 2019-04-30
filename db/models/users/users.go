package users

import (
	"database/sql"
	"time"
)

// User user
type User struct {
	ID         int
	Login      string
	FirstName  string
	LastName   string
	Email      string
	Registered time.Time
	Active     bool
}

// conn DB connection
var conn *sql.DB

// CheckLoginExists checks if login exists
func CheckLoginExists(login string) (bool, error) {
	return true, nil
}

// Create creates new user
func Create(u User) (*User, error) {
	return nil, nil
}

// Update updates user without password
func (u *User) Update() error {
	return nil
}

// UpdatePassword updates password
func (u *User) UpdatePassword(newPassword string) error {
	return nil
}

// Login checks login and password
func Login(login string, password string) (bool, int, error) {
	rows, err := conn.Query("SELECT id from users WHERE login=$1 AND password=$2", login, password)
	if err != nil {
		return false, 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var userID int
		if err = rows.Scan(&userID); err == nil {
			return true, userID, err
		}
	}
	return false, 0, nil
}

// SetConn sets DB connection
func SetConn(c *sql.DB) {
	conn = c
}

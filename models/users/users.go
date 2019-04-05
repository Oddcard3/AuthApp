package users

// User user information
type User struct {
	login   string
	pwdHash string
	pwdSalt string
}

var users = []User{{"user", "123456", ""}}

// Get gets user by login
func Get(login string) (u *User, e error) {
	e = nil
	u = nil
	for _, v := range users {
		if v.login == login {
			u = &v //TODO: copy
			break
		}
	}
	return
}

// CheckPassword checks password
func (u *User) CheckPassword(pwd string) bool {
	return u.pwdHash == pwd
}

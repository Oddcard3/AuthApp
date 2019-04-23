package db

import (
	"authapp/db/models"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq" // postgrsql driver
	"github.com/spf13/viper"
)

var conn *sql.DB

// OpenDB opens DB
func OpenDB() (db *sql.DB, err error) {
	connStr := viper.GetString("db.url")
	log.WithFields(log.Fields{"url": connStr}).Debug("DB connecting...")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to open DB")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to ping DB")
		db = nil
		return nil, err
	}
	conn = db

	models.SetConn(db)
	return conn, nil
}

// Create create DB
func Create() (err error) {
	if conn != nil {
		err = errors.New("No connection to DB")
		log.WithFields(log.Fields{"err": err}).Error("DB conn is nil")
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id serial primary key,
			login text not null UNIQUE,
			email text,
			first_name text,
			second_name text,
			password text not null,
			salt text not null,
			reg_time timestamp not null,
			active boolean not null);`

	_, err = conn.Exec(query)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create users table")
		return
	}

	query = `
		CREATE TABLE IF NOT EXISTS chats (
			id serial primary key,
			creator integer REFERENCES users (id),
			create_ts timestamp not null,
			name text);`

	_, err = conn.Exec(query)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create chats table")
		return
	}

	query = `
		CREATE TABLE IF NOT EXISTS chat_members (
			chat_id integer REFERENCES chats (id),
			user_id integer REFERENCES users (id),
			added_ts timestamp not null,
			PRIMARY KEY(chat_id, user_id));`

	_, err = conn.Exec(query)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create chat_members table")
		return
	}

	query = `
		CREATE TABLE IF NOT EXISTS messages (
			id serial primary key,
			chat_id integer REFERENCES chats (id),
			src_user_id integer REFERENCES users (id),
			ts timestamp not null,
			message text not null);`

	_, err = conn.Exec(query)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create messages table")
		return
	}

	return
}

func createUser(login string, password string) (id int, err error) {
	query :=
		`INSERT INTO users VALUES (
			nextval('users_id_seq'),
			$1,
			'',
			'',
			'',
			$2,
			'',
			NOW(),
			TRUE) RETURNING id`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(login, password).Scan(&id)
	return
}

func createChat(userID int) (id int, err error) {
	query :=
		`INSERT INTO chats VALUES (
			nextval('chats_id_seq'),
			$1,
			NOW(),
			'') RETURNING id`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(userID).Scan(&id)
	return
}

func addUserToChat(chatID int, userID int) (err error) {
	query :=
		`INSERT INTO chat_members VALUES (
			$1,
			$2,
			NOW())`
	_, err = conn.Exec(query, chatID, userID)
	return
}

// FillTestData fills table with test data
func FillTestData() (err error) {
	id1, err := createUser("robocop", "stalone123")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create test users")
		return
	}
	id2, err := createUser("terminator", "arnold123")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create test users")
		return
	}

	chatID, err := createChat(id1)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to create test chat")
		return
	}

	err = addUserToChat(chatID, id1)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to add user to chat")
		return
	}

	err = addUserToChat(chatID, id2)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to add user to chat")
		return
	}
	// query := `
	// 	INSERT INTO users VALUES (
	// 		nextval('users_id_seq'),
	// 		$1,
	// 		'',
	// 		'',
	// 		'',
	// 		'stalone123',
	// 		'',
	// 		NOW(),
	// 		TRUE);
	// 	INSERT INTO users VALUES (
	// 		nextval('users_id_seq'),
	// 		$2,
	// 		'',
	// 		'',
	// 		'',
	// 		'arnold123',
	// 		'',
	// 		NOW(),
	// 		TRUE);`

	// _, err = conn.Exec(query, user1, user2)
	// if err != nil {
	// 	log.WithFields(log.Fields{"err": err}).Error("Failed to create test users")
	// 	return
	// }

	return
}

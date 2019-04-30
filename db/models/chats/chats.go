package chats

import (
	"authapp/db/models/messages"
	"authapp/db/models/users"
	"database/sql"
)

// Chat chat
type Chat struct {
	ID    int
	Users []users.User
	// Creator int
	// Created time.Time
	// Name    string
}

// conn DB connection
var conn *sql.DB

// GetChat gets chat instance by ID
func GetChat(chatID int) (chat *Chat, err error) {
	rows, err := conn.Query("SELECT u.id, u.login FROM chat_members cm, users u WHERE cm.user_id = u.id AND cm.chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chat = new(Chat)
	chat.ID = chatID
	chat.Users = make([]users.User, 0)
	for rows.Next() {
		var login string
		var userID int
		if err = rows.Scan(&userID, &login); err != nil {
			return
		}
		var user users.User
		user.ID = userID
		user.Login = login
		chat.Users = append(chat.Users, user)
	}
	return chat, nil
}

// GetByUser gets chats by user ID
func GetByUser(userID int) (chats []*Chat, err error) {
	rows, err := conn.Query("SELECT chat_id FROM chat_members WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats = make([]*Chat, 0)
	for rows.Next() {
		var chatID int
		if err = rows.Scan(&chatID); err != nil {
			return
		}
		chat, err := GetChat(chatID)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return
}

// GetLastMessages gets last num messages
func GetLastMessages(chatID int, num int) (msgList []messages.Message, err error) {
	rows, err := conn.Query("SELECT id, chat_id, src_user_id, ts, message FROM messages WHERE chat_id=$1 ORDER BY ts DESC LIMIT 100", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgList = make([]messages.Message, 0)
	for rows.Next() {
		var msg messages.Message
		if err = rows.Scan(&msg.ID, &msg.ChatID, &msg.Creator, &msg.Created, &msg.Text); err != nil {
			return
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}

// SearchMessagesForChat searches messages in chat
func SearchMessagesForChat(chatID int, text string) ([]messages.Message, error) {
	return nil, nil
}

// SearchMessagesForUser searches messages in all user chats
func SearchMessagesForUser(userID int, text string) ([]messages.Message, error) {
	return nil, nil
}

// AddMessage adds message to chat
func AddMessage(msg *messages.Message) error {
	query :=
		`INSERT INTO messages VALUES (
		$1,
		$2,
		$3,
		NOW(),
		$4)`
	_, err := conn.Exec(query, msg.ID, msg.ChatID, msg.Creator, msg.Text)
	return err
}

// Create creates new chat
func Create(userID int, Users []int) (*Chat, error) {
	return nil, nil
}

// SetConn sets DB connection
func SetConn(c *sql.DB) {
	conn = c
}

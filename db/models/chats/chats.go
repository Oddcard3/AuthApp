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

// GetMembers gets chat members
func GetMembers(chatID int) (members []users.User, err error) {
	rows, err := conn.Query("SELECT u.id, u.login FROM chat_members cm, users u WHERE cm.user_id = u.id AND cm.chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	members = make([]users.User, 0)
	for rows.Next() {
		var login string
		var userID int
		if err = rows.Scan(&userID, &login); err != nil {
			return
		}
		var user users.User
		user.ID = userID
		user.Login = login
		members = append(members, user)
	}
	return members, nil
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
		chat := new(Chat)
		chat.ID = chatID
		chat.Users = make([]users.User, 0)
		chats = append(chats, chat)
	}

	for _, c := range chats {
		c.Users, err = GetMembers(c.ID)
		if err != nil {
			return
		}
	}

	return
}

// GetLastMessages gets last num messages
func GetLastMessages(chatID int, num int) (msgList []messages.Message, err error) {
	rows, err := conn.Query("SELECT id, chat_id, src_user_id, ts, message  FROM messages chat_id = $1", chatID)
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
func AddMessage(chatID int, msg messages.Message) (*messages.Message, error) {
	return nil, nil
}

// Create creates new chat
func Create(userID int, Users []int) (*Chat, error) {
	return nil, nil
}

// SetConn sets DB connection
func SetConn(c *sql.DB) {
	conn = c
}

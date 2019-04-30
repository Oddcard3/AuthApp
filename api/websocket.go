package api

import (
	"container/list"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	jwttoken "authapp/auth/jwt"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"authapp/db/models/chats"
	"authapp/db/models/messages"
)

//var upgrader = websocket.Upgrader{} // use default options

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsMsg msg
type WsMsg struct {
	Type string `json:"event"`
}

// WsUser user
type WsUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// WsChatMessage message
type WsChatMessage struct {
	ID     string    `json:"id"`
	ChatID int       `json:"chatId"`
	Ts     time.Time `json:"ts"`
	UserID int       `json:"userId"`
	Text   string    `json:"text"`
}

// WsChat chat
type WsChat struct {
	ID       int			 `json:"id"`
	Members  []WsUser        `json:"users"`
	Messages []WsChatMessage `json:"messages"`
}

// WsChats chats
type WsChats struct {
	UserID int `json:"userId"`
	Chats []WsChat `json:"chats"`
}

// WsChatsEvent chats event
type WsChatsEvent struct {
	Event string   `json:"event"`
	Data  *WsChats `json:"data"`
}

// WsMsgEvent chats event
type WsMsgEvent struct {
	Event string         `json:"event"`
	Msg   *WsChatMessage `json:"data"`
}

var connListLock = new(sync.RWMutex)
var connList = list.New()

type clientConn struct {
	// websocket.Conn is not thread safe (https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency)
	ConnLock sync.Mutex
	Conn     *websocket.Conn
	UserID   int
}

func addConnection(userID int, conn *websocket.Conn) (*clientConn, error) {
	connListLock.Lock()
	defer connListLock.Unlock()

	clConn := new(clientConn)
	clConn.Conn = conn
	clConn.UserID = userID

	connList.PushBack(clConn)
	return clConn, nil
}

func removeConnection(conn *websocket.Conn) error {
	connListLock.Lock()
	defer connListLock.Unlock()

	var next *list.Element
	for e := connList.Front(); e != nil; e = next {
		clConn := e.Value.(*clientConn)
		next = e.Next()
		if clConn.Conn == conn {
			connList.Remove(e)
		}
	}
	return nil
}

func isUserInChat(userID int, chat *chats.Chat) bool {
	for _, v := range chat.Users {
		if v.ID == userID {
			return true
		}
	}
	return false
}

func sendMessage(msgEvent *WsMsgEvent) error {
	msg := msgEvent.Msg
	chat, err := chats.GetChat(msg.ChatID)
	if err != nil {
		return err
	}

	msgJSON, err := json.Marshal(msgEvent)
	if err != nil {
		return err
	}

	connListLock.RLock()
	defer connListLock.RUnlock()
	for e := connList.Front(); e != nil; e = e.Next() {
		clConn := e.Value.(*clientConn)

		if isUserInChat(clConn.UserID, chat) {
			clConn.ConnLock.Lock()
			defer clConn.ConnLock.Unlock()

			if err := clConn.Conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
				log.WithFields(log.Fields{"err": err, "chatId": chat.ID,
					"userId": clConn.UserID, "msg": string(msgJSON)}).Error("Failed to send message to user")
				return err
			}
			log.WithFields(log.Fields{"chatId": chat.ID, "userId": clConn.UserID, "msg": string(msgJSON)}).Info("Msg sent to user")
		}
	}
	return nil
}

func sendChats(clConn *clientConn) error {
	chatsList, err := chats.GetByUser(clConn.UserID)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "userID": clConn.UserID}).Error("Failed to get chats for user")
		return err
	}

	wsChats := make([]WsChat, len(chatsList))
	for i, chat := range chatsList {
		wsChats[i].ID = chat.ID
		wsChats[i].Members = make([]WsUser, len(chat.Users))
		for j, u := range chat.Users {
			wsChats[i].Members[j].ID = u.ID
			wsChats[i].Members[j].Login = u.Login
		}



		msgs, err := chats.GetLastMessages(chat.ID, 100)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "chatID": chat.ID}).Error("Failed to get last messages for chat")
			return err
		}

		wsChats[i].Messages = make([]WsChatMessage, len(msgs))
		for j, msg := range msgs {
			wsChats[i].Messages[j].ID = msg.ID
			wsChats[i].Messages[j].ChatID = msg.ChatID
			wsChats[i].Messages[j].Ts = msg.Created
			wsChats[i].Messages[j].UserID = msg.Creator
			wsChats[i].Messages[j].Text = msg.Text
		}
	}

	event := &WsChatsEvent{"chats", &WsChats{clConn.UserID, wsChats}}

	clConn.ConnLock.Lock()
	defer clConn.ConnLock.Unlock()

	chatsJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	clConn.Conn.WriteMessage(websocket.TextMessage, chatsJSON)
	return nil
}

func storeMessage(msg *WsChatMessage) error {
	dbMsg := &messages.Message{msg.ID, msg.Text, msg.Ts, msg.ChatID, msg.UserID}
	return chats.AddMessage(dbMsg)
}

func handleNewMessageEvent(userID int, data []byte, conn *websocket.Conn) {
	msgEvent := new(WsMsgEvent)
	if err := json.Unmarshal(data, msgEvent); err != nil {
		log.WithFields(log.Fields{"err": err, "data": string(data)}).Error("Failed to parse WsMsgEvent")
		return
	}

	// sets current user id
	msgEvent.Msg.UserID = userID

	if err := sendMessage(msgEvent); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to send message to chat")
	}

	if err := storeMessage(msgEvent.Msg); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to store message")
	}
}

func handleEvent(userID int, data []byte, conn *websocket.Conn) error {

	unquotedData, _ := strconv.Unquote(string(data))

	var msg WsMsg
	if err := json.Unmarshal([]byte(unquotedData), &msg); err != nil {
		log.WithFields(log.Fields{"err": err, "data": unquotedData}).Error("Failed to parse message")
		return err
	}

	switch msg.Type {
	case "send-msg":
		handleNewMessageEvent(userID, []byte(unquotedData), conn)
	default:
		log.WithFields(log.Fields{"type": msg.Type}).Error("Incorrect message type received")
	}
	return nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")

	userIDStr, err := jwttoken.GetUserID(tokenStr)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "token": tokenStr}).Error("No user for token")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "userID": userIDStr}).Error("Incorrect user ID, must be integer")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("upgrade error:")
		return
	}
	defer c.Close()

	clConn, err := addConnection(userID, c)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to add new connection")
		return
	}
	if err := sendChats(clConn); err != nil {
		log.WithFields(log.Fields{"err": err, "userID": clConn.UserID}).Error("Failed to send chats")
		return
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("read error:")
			break
		}
		log.WithFields(log.Fields{"msg": message}).Info("recv: ")

		handleEvent(userID, message, c)
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.WithFields(log.Fields{"err": err}).Error("write error:")
		// 	break
		// }
	}
	removeConnection(c)
}

// WebsocketRouter handler for websocket
func WebsocketRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		//r.Use(jwttoken.GetVerifier())
		//r.Use(jwttoken.AppAuthenticator)
		r.Get("/{token}", handleConnection)
	})
	return r
}

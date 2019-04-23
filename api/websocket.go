package api

import (
	"net/http"
	"time"

	// "github.com/go-chi/jwtauth"

	jwttoken "authapp/auth/jwt"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	// "authapp/db/models/messages"
	// "authapp/db/models/users"
	// "authapp/db/models/chats"
)

var upgrader = websocket.Upgrader{} // use default options

// WsMsg msg
type WsMsg struct {
	MsgType string `json:"type"`
}

// WsUser user
type WsUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// WsChatMessage message
type WsChatMessage struct {
	ID     int       `json:"id"`
	ChatID int       `json:"chat_id"`
	Ts     time.Time `json:"ts"`
	UserID int       `json:"user_id"`
	Text   string    `json:"text"`
}

// WsChat chat
type WsChat struct {
	ID       int
	Members  []WsUser        `json:"users"`
	Messages []WsChatMessage `json:"messages"`
}

// WsChats chats
type WsChats struct {
	Chats []WsChat `json:"chats"`
}

// WsChatsEvent chats event
type WsChatsEvent struct {
	Event string  `json:"event"`
	Data  WsChats `json:"data"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	// userID, err := jwttoken.GetUserID()
	// if err != nil {
	// 	http.Error(w, http.StatusText(401), 401)
	//  return
	// }

	// _, claims, _ := jwtauth.FromContext(r.Context())
	// userID, ok := claims["user_id"]
	// log.WithFields(log.Fields{"user_id": userID}).Error("User ID connection")
	// if !ok {
	// 	http.Error(w, http.StatusText(401), 401)
	// 	return
	// }

	// secHeader := r.Header.Get("Sec-WebSocket-Protocol")
	// s := strings.Split(secHeader, ",")
	// if len(s) < 2 {
	// 	log.WithFields(log.Fields{"header": secHeader}).Error("Incorrect format of Sec-WebSocket-Protocol")
	// 	http.Error(w, http.StatusText(401), 401)
	// 	return
	// }

	// tokenStr := strings.TrimSpace(s[1])

	_, err := jwttoken.GetUserID(tokenStr)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "token": tokenStr}).Error("No user for token")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("upgrade error:")
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("read error:")
			break
		}
		log.WithFields(log.Fields{"msg": message}).Info("recv: ")
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("write error:")
			break
		}
	}
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

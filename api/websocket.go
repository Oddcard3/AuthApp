package api

import (
	"net/http"

	. "authapp/logging"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Logger.WithFields(logrus.Fields{"err": err}).Error("upgrade error:")
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			Logger.WithFields(logrus.Fields{"err": err}).Error("read error:")
			break
		}
		Logger.WithFields(logrus.Fields{"msg": message}).Info("recv: ")
		err = c.WriteMessage(mt, message)
		if err != nil {
			Logger.WithFields(logrus.Fields{"err": err}).Error("write error:")
			break
		}
	}
}

// WebsocketRouter handler for websocket
func WebsocketRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handleConnection)

	return r
}

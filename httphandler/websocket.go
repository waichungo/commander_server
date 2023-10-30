package httphandler

import (
	"commander_server/db"
	"commander_server/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gorm.io/gorm"
)

type WebsocketMessage struct {
	To   string `json:"to"`
	From string `json:"from"`
	Data string `json:"data"`
	Type int    `json:"type"`
}

type MessageType int

const (
	UNDEFINED MessageType = iota
	CLIENTNOTFOUND
	SERVERERROR
)

var _clients map[string]*websocket.Conn
var clientsLck = sync.Mutex{}

func ExecuteOnWebSocketClients(execute func(clients map[string]*websocket.Conn)) {
	clientsLck.Lock()
	defer clientsLck.Unlock()
	execute(_clients)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ValidateClient(clientId string) bool {
	if len(clientId) > 0 {
		err := db.ExecuteOnDB(func(db *gorm.DB) error {
			cl := models.Client{}
			err := db.Model(&cl).Find(&cl).Error
			return err
		})
		return err == nil
	}
	return false
}
func HandleWebsocket(c *gin.Context) {
	clientId := strings.TrimSpace(c.GetHeader("X-CLIENT-ID"))
	if len(clientId) > 0 {
		if _, f := _clients[clientId]; f {
			c.AbortWithError(http.StatusBadRequest, errors.New("client already exists"))
			return
		}
	}
	if ValidateClient(clientId) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// panic(err)
			log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ExecuteOnWebSocketClients(func(clients map[string]*websocket.Conn) {
			clients[clientId] = conn
		})
		for {
			// Read message from client
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				// panic(err)
				log.Printf("%s, error while reading message\n", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				break
			}
			if messageType == websocket.TextMessage {
				msg := WebsocketMessage{}
				err = json.Unmarshal(p, &msg)
				if err == nil {
					socket, f := _clients[msg.To]
					if f {
						socket.WriteMessage(websocket.TextMessage, p)
					} else {
						conn.WriteJSON(&WebsocketMessage{
							From: "server",
							To:   msg.From,
							Type: int(CLIENTNOTFOUND),
						})

					}
				} else {
					conn.WriteJSON(&WebsocketMessage{
						From: "server",
						To:   msg.From,
						Type: int(SERVERERROR),
						Data: "Error parsing message.A valid type is required",
					})
				}
			}

		}
		ExecuteOnWebSocketClients(func(clients map[string]*websocket.Conn) {
			delete(clients, clientId)
		})
	} else {
		c.AbortWithError(http.StatusBadRequest, errors.New("a valid client Id is required"))
	}

}

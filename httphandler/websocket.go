package httphandler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	To   string      `json:"to"`
	From string      `json:"from"`
	Data string      `json:"data"`
	Type MessageType `json:"type"`
}

type MessageType int

const (
	UNDEFINED MessageType = iota
	CLIENTNOTFOUND
	SERVERERROR
	WEBSOCKETCONNECTIONLIST
)

type WebsocketInfo struct {
	Conn      *websocket.Conn
	IsMachine bool
}

var _clients map[string]WebsocketInfo = map[string]WebsocketInfo{}
var clientsLck = sync.Mutex{}

func ExecuteOnWebSocketClients(execute func(clients map[string]WebsocketInfo)) {
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

func GetIdentity(c *gin.Context) *Identifier {
	defer func() {
		err := recover()
		if err != nil {
			println(err)
		}
	}()
	identifier := Identifier{}

	claims := jwt.ExtractClaims(c)

	identifier.Identifier = claims["identifier"].(string)
	identifier.Name = claims["name"].(string)
	identifier.Type = int(claims["type"].(float64))
	return &identifier

}

//	func ValidateClient(clientId string) bool {
//		if len(clientId) > 0 {
//			var err error
//			db.ExecuteOnDB(func(db *gorm.DB) error {
//				cl := models.Client{}
//				err = db.Model(&cl).Where("id = ?").First(&cl).Error
//				return err
//			})
//			return err == nil
//		}
//		return false
//	}
type RealtimeClient struct {
	Id        string `json:"id"`
	IsMachine bool   `json:"isMachine"`
}

func GetIdentityFromHeader(c *gin.Context) (*Identifier, error) {
	var identity *Identifier
	var err error
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		err = errors.New("unauthorized")
		return nil, err
	}
	if !strings.Contains(auth, "Bearer") {
		err = errors.New("unauthorized")
		return nil, err
	}
	authSplit := strings.Split(auth, " ")
	if len(authSplit) != 2 {
		err = errors.New("invalid authorization")
		return nil, err
	}
	takenSplit := strings.Split(authSplit[1], ".")
	if len(takenSplit) != 3 {
		err = errors.New("invalid authorization")
		return nil, err
	}
	var dec []byte
	dec, err = base64.URLEncoding.DecodeString(takenSplit[1])
	if err != nil {
		err = errors.New("invalid authorization")
		return nil, err
	}
	id := Identifier{}
	err = json.Unmarshal(dec, &id)

	if err != nil {
		err = errors.New("invalid authorization")
		return nil, err
	} else {
		identity = &id
	}
	return identity, err
}
func GetLiveClients() []RealtimeClient {
	var res []RealtimeClient
	ExecuteOnWebSocketClients(func(clients map[string]WebsocketInfo) {
		res = make([]RealtimeClient, 0, len(clients))
		for key, v := range clients {
			res = append(res, RealtimeClient{
				Id:        key,
				IsMachine: v.IsMachine,
			})
		}
	})
	return res
}
func init() {
	go func() {
		count := 0
		for {

			time.Sleep(time.Second)
			ExecuteOnWebSocketClients(func(clients map[string]WebsocketInfo) {
				if len(clients) > 0 {
					count++
					msg := fmt.Sprintf("Counter %d", count)
					for _, cl := range clients {
						cl.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
					}

				} else {
					count = 0
				}

			})
		}
	}()
}
func HandleWebsocket(c *gin.Context) {
	identity, err := GetIdentityFromHeader(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	// if identity != nil {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{})
	if err != nil {
		// panic(err)
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ExecuteOnWebSocketClients(func(clients map[string]WebsocketInfo) {
		clients[identity.Identifier] = WebsocketInfo{
			Conn:      conn,
			IsMachine: identity.Type == 2,
		}
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
					socket.Conn.WriteMessage(websocket.TextMessage, p)
				} else {
					if msg.Type == WEBSOCKETCONNECTIONLIST {
						cl := GetLiveClients()
						data, _ := json.Marshal(cl)
						conn.WriteJSON(&WebsocketMessage{
							From: "server",
							To:   msg.From,
							Type: CLIENTNOTFOUND,
							Data: string(data),
						})

					} else {
						conn.WriteJSON(&WebsocketMessage{
							From: "server",
							To:   msg.From,
							Type: CLIENTNOTFOUND,
						})
					}

				}
			} else {
				conn.WriteJSON(&WebsocketMessage{
					From: "server",
					To:   msg.From,
					Type: SERVERERROR,
					Data: "Error parsing message.A valid type is required",
				})
			}
		}

	}
	ExecuteOnWebSocketClients(func(clients map[string]WebsocketInfo) {
		delete(clients, identity.Identifier)
	})
	// } else {
	// 	c.AbortWithError(http.StatusBadRequest, errors.New("a valid client Id is required"))
	// }

}

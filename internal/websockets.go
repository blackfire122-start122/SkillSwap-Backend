package internal

import (
	"SkillSwap/pkg"
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Client struct {
	Conn   *websocket.Conn
	RoomId string
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

var clients = make(map[Client]bool)

func handleConnections(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	roomId := c.Param("roomId")

	if roomId == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Println("roomId empty")
		return
	}

	var chat pkg.SkillChat

	if err := pkg.DB.Where("id=?", roomId).First(&chat).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chat.CustomerID != user.Id && chat.PerformerID != user.Id {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	client := Client{Conn: conn, RoomId: roomId}
	clients[client] = true

	for {
		var msg ClientMessage

		err := conn.ReadJSON(&msg)

		if err != nil {
			fmt.Println("readJson ", err)
			delete(clients, client)
			return
		}

		if msg.Type == "msg" {
			content, ok := msg.Content.(string)
			if !ok {
				fmt.Println("msg content not string")
				delete(clients, client)
				return
			}

			message := pkg.RedisMessage{CreatedAt: time.Now(), Message: content, UserId: user.Id, ChatId: uint64(chat.ID)}
			messageJSON, err := json.Marshal(message)
			if err != nil {
				fmt.Println("json marshal error", err)
				delete(clients, client)
				return
			}

			indexMsg, err := pkg.RedisClient.LPush(pkg.Ctx, "chat_messages", messageJSON).Result()
			if err != nil {
				fmt.Println("redis lpush error", err)
				dbMessage := pkg.Message{Message: content, Read: false, User: user, SkillChatID: message.ChatId}
				if err := pkg.DB.Create(&dbMessage).Error; err != nil {
					fmt.Println("db err save message ", err)
					delete(clients, client)
					return
				}
			}

			msgResp := make(map[string]interface{})
			msgResp["message"] = message.Message
			msgResp["userId"] = user.Id
			msgResp["createdAt"] = message.CreatedAt
			msgResp["indexRedis"] = indexMsg

			err = sendInGroup(ClientMessage{Type: "msg", Content: msgResp}, client)

			if err != nil {
				fmt.Println(err)
				delete(clients, client)
				return
			}
		} else if msg.Type == "changeStatus" {
			err = sendInGroup(msg, client)
			if err != nil {
				fmt.Println(err)
				delete(clients, client)
				return
			}
		} else if msg.Type == "msgRead" {
			pkg.RedisClient.LPush(pkg.Ctx, "chat_messages_readable", msg.Content)
			err = sendInGroup(msg, client)
			if err != nil {
				fmt.Println(err)
				delete(clients, client)
				return
			}
		}
	}
}

func sendInGroup(msg interface{}, client Client) error {
	for cl := range clients {
		if cl.RoomId == client.RoomId {
			err := cl.Conn.WriteJSON(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SendPing() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for client := range clients {
			if err := client.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				fmt.Println("error send ping", err)
				break
			}
		}
	}
}

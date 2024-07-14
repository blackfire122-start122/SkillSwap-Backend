package internal

import (
	"SkillSwap/pkg"
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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

			err = pkg.RedisClient.LPush(pkg.Ctx, "chat_messages", messageJSON).Err()
			if err != nil {
				fmt.Println("redis lpush error", err)
				delete(clients, client)
				return
			}

			msgResp := make(map[string]interface{})
			msgResp["message"] = message.Message
			msgResp["userId"] = user.Id
			msgResp["createdAt"] = message.CreatedAt

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

func SaveMessagesToDB() {
	const batchSize = 100
	for {
		messages, err := pkg.RedisClient.LRange(pkg.Ctx, "chat_messages", 0, batchSize-1).Result()
		if err != nil {
			fmt.Println("redis lrange error", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		if len(messages) == 0 {
			time.Sleep(5 * time.Minute)
			continue
		}

		if err := pkg.RedisClient.LTrim(pkg.Ctx, "chat_messages", int64(batchSize), -1).Err(); err != nil {
			fmt.Println("redis ltrim error", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		var redisMessages []pkg.RedisMessage
		for _, messageJSON := range messages {
			var redisMessage pkg.RedisMessage
			if err := json.Unmarshal([]byte(messageJSON), &redisMessage); err != nil {
				fmt.Println("json unmarshal error", err)
				continue
			}
			redisMessages = append(redisMessages, redisMessage)
		}

		tx := pkg.DB.Begin()
		if tx.Error != nil {
			fmt.Println("db transaction begin error", tx.Error)
			continue
		}

		for _, redisMessage := range redisMessages {
			message := pkg.Message{
				Model:   gorm.Model{CreatedAt: redisMessage.CreatedAt},
				Message: redisMessage.Message,
				Read:    false,
				UserID:  redisMessage.UserId,
			}

			if err := tx.Create(&message).Error; err != nil {
				fmt.Println("db create error", err)
				tx.Rollback()
				break
			}

			var chat pkg.SkillChat
			if err := tx.First(&chat, redisMessage.ChatId).Error; err != nil {
				fmt.Println("db find chat error", err)
				tx.Rollback()
				break
			}

			if err := tx.Model(&chat).Association("Messages").Append(&message); err != nil {
				fmt.Println("db append message error", err)
				tx.Rollback()
				break
			}
		}

		if err := tx.Commit().Error; err != nil {
			fmt.Println("db transaction commit error", err)
		}
	}
}

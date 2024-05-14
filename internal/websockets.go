package internal

import (
	//. "gameServer/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Message struct {
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

	//var user User

	//DB.First(&user)

	client := Client{Conn: conn, RoomId: roomId}
	clients[client] = true

	for {
		var msg Message

		err := conn.ReadJSON(&msg)

		if err != nil {
			fmt.Println("readJson ", err)
			delete(clients, client)
			return
		}

		fmt.Println(msg)

		//if msg.Type == "newPlayer" {
		//	err = sendInGroup(msg, client)
		//	if err != nil {
		//		fmt.Println(err)
		//		delete(clients, client)
		//		return
		//	}
		//}
		//if msg.Type == "move" {
		//	err = sendInGroup(msg, client)
		//	if err != nil {
		//		fmt.Println(err)
		//		delete(clients, client)
		//		return
		//	}
		//}
	}
}

//func sendInGroup(msg Message, client Client) error {
//	for cl := range clients {
//		if cl.RoomId == client.RoomId && cl != client {
//			err := cl.Conn.WriteJSON(msg)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}

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

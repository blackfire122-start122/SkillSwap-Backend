package pkg

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func SkillChatMessages(chat SkillChat) ([]map[string]interface{}, error) {
	resp := make([]map[string]interface{}, 0)
	for _, message := range chat.Messages {
		item := make(map[string]interface{})
		item["id"] = message.ID
		item["message"] = message.Message
		item["userId"] = message.UserID
		item["createdAt"] = message.CreatedAt
		resp = append(resp, item)
	}
	return resp, nil
}

func GetSkillChatData(chatId uint64, user User) (map[string]interface{}, error) {
	var skillChat SkillChat
	resp := make(map[string]interface{})

	if err := DB.Preload("Status").First(&skillChat, chatId).Error; err != nil {
		return resp, err
	}

	if skillChat.CustomerID != user.Id && skillChat.PerformerID != user.Id {
		return resp, errors.New("Forbidden")
	}

	resp["id"] = skillChat.ID
	resp["status"] = skillChat.Status.Status
	resp["performerID"] = skillChat.PerformerID
	resp["customerID"] = skillChat.CustomerID

	return resp, nil
}

func GetChatRedisMessages(chatId uint64) []Message {
	redisMessages, err := RedisClient.LRange(Ctx, "chat_messages", 0, -1).Result()
	if err != nil {
		fmt.Println("redis lrange error", err)
	}

	var messagesFromRedis []Message
	for _, messageJSON := range redisMessages {
		var redisMessage RedisMessage
		if err := json.Unmarshal([]byte(messageJSON), &redisMessage); err != nil {
			fmt.Println("json unmarshal error", err)
			continue
		}

		if redisMessage.ChatId == chatId {
			messagesFromRedis = append(messagesFromRedis, Message{
				Model:       gorm.Model{CreatedAt: redisMessage.CreatedAt},
				Message:     redisMessage.Message,
				Read:        false,
				UserID:      redisMessage.UserId,
				SkillChatID: redisMessage.ChatId,
			})
		}
	}

	return messagesFromRedis
}

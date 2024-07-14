package controllers

import (
	"SkillSwap/pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Order(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var orderSkill pkg.OrderSkill
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &orderSkill); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := pkg.CreateChat(user, orderSkill)

	if err != nil {
		if err.Error() == "chat already exists" {
			c.JSON(http.StatusAlreadyReported, resp)
			return
		}
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetCustomerSkillChats(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := pkg.DB.Preload("CustomerSkillChats").Find(&user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := pkg.GetChatsCustomerData(user.CustomerSkillChats)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetPerformerSkillChats(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := pkg.DB.Preload("PerformerSkillChats").Find(&user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := pkg.GetChatsPerformerData(user.PerformerSkillChats)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetSkillChatMessages(c *gin.Context) {
	const MAX_MESSAGES_FROM_DB = 20
	const MAX_MESSAGES_FROM_REDIS_WITH_DB = 100

	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	countMessages, err := strconv.ParseUint(c.Query("countMessages"), 10, 64)

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseUint(c.Query("chatId"), 10, 64)

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	messagesFromRedis := pkg.GetChatRedisMessages(chatId)
	countMessagesFromRedis := len(messagesFromRedis)
	needRedisMessages := countMessagesFromRedis - int(countMessages)

	if needRedisMessages > 0 {
		messagesFromRedis = messagesFromRedis[countMessagesFromRedis-needRedisMessages:]
	} else {
		messagesFromRedis = []pkg.Message{}
	}

	var skillChat pkg.SkillChat

	if err := pkg.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Offset(int(countMessages) - countMessagesFromRedis).Limit(MAX_MESSAGES_FROM_DB)
	}).First(&skillChat, chatId).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if skillChat.CustomerID != user.Id && skillChat.PerformerID != user.Id {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	allMessages := append(skillChat.Messages, messagesFromRedis...)

	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].CreatedAt.After(allMessages[j].CreatedAt)
	})

	end := MAX_MESSAGES_FROM_REDIS_WITH_DB
	if end > len(allMessages) {
		end = len(allMessages)
	}

	messages := allMessages[:end]

	skillChat.Messages = messages

	resp, err := pkg.SkillChatMessages(skillChat)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetChat(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	chatId, err := strconv.ParseUint(c.Query("chatId"), 10, 64)

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := pkg.GetSkillChatData(chatId, user)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetStatuses(c *gin.Context) {
	loginUser, _ := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := pkg.GetStatusesData()

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func SetStatus(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var putStatusRequest pkg.PutStatusRequest
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &putStatusRequest); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := pkg.SetStatusChat(user, putStatusRequest); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

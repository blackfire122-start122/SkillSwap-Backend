package controllers

import (
	"SkillSwap/pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	var skillChat pkg.SkillChat

	if err := pkg.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Offset(int(countMessages)).Limit(20)
	}).First(&skillChat, c.Query("chatId")).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if skillChat.CustomerID != user.Id && skillChat.PerformerID != user.Id {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	resp, err := pkg.SkillChatMessages(skillChat)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

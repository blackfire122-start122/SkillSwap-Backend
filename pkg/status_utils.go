package pkg

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var Canceled string = "Canceled"
var DiscussionOfDetails string = "Discussion of details"
var Executing string = "Executing"
var Executed string = "Executed"
var Paid string = "Paid"
var Rewarded string = "Rewarded"

var StatusTexts = []string{
	Canceled,
	DiscussionOfDetails,
	Executing,
	Executed,
	Paid,
	Rewarded,
}

func CreateDefaultModelData() {
	var statuses []*Status

	for _, statusText := range StatusTexts {
		statuses = CheckAndCreateStatus(statusText, statuses)
	}

	if len(statuses) > 0 {
		if err := DB.Create(&statuses).Error; err != nil {
			fmt.Print(err)
		}
	}
}

func CheckAndCreateStatus(StatusText string, statuses []*Status) []*Status {
	var status Status
	if err := DB.First(&status).Where("status=?", StatusText).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			statuses = append(statuses, &Status{Status: StatusText})
		}
	}

	return statuses
}

func GetStatusesData() ([]map[string]interface{}, error) {
	var statuses []Status
	resp := make([]map[string]interface{}, 0)

	if err := DB.Find(&statuses).Error; err != nil {
		return resp, err
	}

	for _, status := range statuses {
		item := make(map[string]interface{})
		item["id"] = status.ID
		item["status"] = status.Status
		resp = append(resp, item)
	}
	return resp, nil
}

type PutStatusRequest struct {
	StatusId uint64 `json:"statusId"`
	ChatId   uint64 `json:"chatId"`
}

func SetStatusChat(user User, putStatusRequest PutStatusRequest) error {
	var status Status
	var skillChat SkillChat

	if err := DB.First(&status, putStatusRequest.StatusId).Error; err != nil {
		return err
	}

	if err := DB.First(&skillChat, putStatusRequest.ChatId).Error; err != nil {
		return err
	}

	if skillChat.CustomerID != user.Id && skillChat.PerformerID != user.Id {
		return errors.New("Forbidden")
	}

	skillChat.Status = status

	if err := DB.Save(&skillChat).Error; err != nil {
		return err
	}

	return nil
}

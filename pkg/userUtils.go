package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func GetUserData(user User) (map[string]interface{}, error) {
	resp := make(map[string]interface{})

	if err := DB.Preload("Skills").Preload("Categories").Find(&user).Error; err != nil {
		return resp, err
	}

	skills := make([]map[string]string, 0)

	for _, skill := range user.Skills {
		item := make(map[string]string)

		item["id"] = strconv.FormatUint(skill.Id, 10)
		item["name"] = skill.Name

		skills = append(skills, item)
	}

	categories := make([]map[string]interface{}, 0)

	for _, category := range user.Categories {
		item := make(map[string]interface{})

		item["id"] = strconv.FormatUint(category.Id, 10)
		item["name"] = category.Name

		if err := DB.Preload("Skills").Find(&category).Error; err != nil {
			return resp, err
		}

		categorySkills := make([]map[string]string, 0)

		for _, skill := range category.Skills {
			itemSkill := make(map[string]string)

			itemSkill["id"] = strconv.FormatUint(skill.Id, 10)
			itemSkill["name"] = skill.Name

			categorySkills = append(categorySkills, itemSkill)
		}

		item["skills"] = categorySkills

		categories = append(categories, item)
	}

	resp["id"] = strconv.FormatUint(user.Id, 10)
	resp["username"] = user.Username
	resp["image"] = user.Image
	resp["email"] = user.Email
	resp["phone"] = user.Phone
	resp["rating"] = user.Rating
	resp["skills"] = skills
	resp["categories"] = categories

	return resp, nil
}

type ReviewData struct {
	Review string `json:"review"`
	Rating uint8  `json:"rating"`
	ToUser uint64 `json:"toUser"`
}

func CreateReviewUser(user User, reviewData ReviewData) (map[string]string, error) {
	resp := make(map[string]string)

	if reviewData.Rating > 100 || reviewData.Rating < 0 {
		return resp, errors.New("rating should be no more than 100 and no less than 0")
	}

	var toUser User

	if err := DB.Preload("Reviews").First(&toUser, reviewData.ToUser).Error; err != nil {
		return resp, err
	}

	var review Review

	if err := DB.Where("reviewer_id = ? AND user_id = ?", user.Id, toUser.Id).First(&review).Error; err == nil {
		review.Rating = reviewData.Rating
		review.Review = reviewData.Review
		if err := DB.Save(&review).Error; err != nil {
			return resp, err
		}
		if err := DB.Model(&toUser).Association("Reviews").Replace(&review); err != nil {
			return resp, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		review = Review{Reviewer: user, Review: reviewData.Review, Rating: reviewData.Rating}

		if err := DB.Create(&review).Error; err != nil {
			return resp, err
		}

		if err := DB.Model(&toUser).Association("Reviews").Append(&review); err != nil {
			return resp, err
		}
	} else {
		return resp, err
	}

	var totalRating uint64
	var reviewCount uint64

	for _, r := range toUser.Reviews {
		totalRating += uint64(r.Rating)
		reviewCount++
	}

	if reviewCount > 0 {
		toUser.Rating = uint8(totalRating / reviewCount)
	} else {
		toUser.Rating = 0
	}

	if err := DB.Save(&toUser).Error; err != nil {
		return resp, err
	}

	resp["Created"] = "OK"
	resp["NewRating"] = strconv.Itoa(int(toUser.Rating))
	resp["Id"] = strconv.Itoa(int(review.Id))
	return resp, nil
}

type OrderSkill struct {
	SkillId uint64 `json:"skillId"`
	ToUser  uint64 `json:"toUser"`
}

func CreateChat(user User, orderSkill OrderSkill) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	var toUser User
	if err := DB.Preload("Skills").Preload("PerformerSkillChats").First(&toUser, orderSkill.ToUser).Error; err != nil {
		return resp, err
	}

	for _, chat := range toUser.PerformerSkillChats {
		if chat.CustomerID == user.Id && chat.SkillID == orderSkill.SkillId {
			resp["chatId"] = chat.ID
			return resp, errors.New("chat already exists")
		}
	}

	var skillFoundFlag = false
	var skillFound Skill

	for _, skill := range toUser.Skills {
		if skill.Id == orderSkill.SkillId {
			skillFoundFlag = true
			skillFound = skill
			break
		}
	}

	if !skillFoundFlag {
		return resp, errors.New("skill not found")
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		var status Status
		if err := DB.First(&status).Where("status=?", DiscussionOfDetails).Error; err != nil {
			fmt.Println(err)
		}

		skillChat := SkillChat{Skill: skillFound, Status: status}

		if err := tx.Create(&skillChat).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("CustomerSkillChats").Append(&skillChat); err != nil {
			return err
		}

		if err := tx.Model(&toUser).Association("PerformerSkillChats").Append(&skillChat); err != nil {
			return err
		}

		resp["chatId"] = skillChat.ID
		return nil
	})

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func GetChatsCustomerData(chats []SkillChat) ([]map[string]interface{}, error) {
	resp := make([]map[string]interface{}, 0)
	for _, chat := range chats {
		if err := DB.Preload("Skill").First(&chat).Error; err != nil {
			return resp, err
		}

		item := make(map[string]interface{})
		item["id"] = chat.ID

		skillItem := make(map[string]interface{})
		skillItem["id"] = chat.Skill.Id
		skillItem["name"] = chat.Skill.Name

		item["skill"] = skillItem

		var performer User
		if err := DB.First(&performer, chat.PerformerID).Error; err != nil {
			return resp, err
		}

		item["performer"] = performer

		resp = append(resp, item)
	}
	return resp, nil
}

func GetChatsPerformerData(chats []SkillChat) ([]map[string]interface{}, error) {
	resp := make([]map[string]interface{}, 0)
	for _, chat := range chats {
		if err := DB.Preload("Skill").First(&chat).Error; err != nil {
			return resp, err
		}

		item := make(map[string]interface{})
		item["id"] = chat.ID

		skillItem := make(map[string]interface{})
		skillItem["id"] = chat.Skill.Id
		skillItem["name"] = chat.Skill.Name

		item["skill"] = skillItem

		var customer User
		if err := DB.First(&customer, chat.CustomerID).Error; err != nil {
			return resp, err
		}

		item["customer"] = customer

		resp = append(resp, item)
	}
	return resp, nil
}

func GenerateJsonObjectUsers(users []User) []map[string]string {
	resp := make([]map[string]string, 0)

	for _, user := range users {
		item := make(map[string]string)

		item["id"] = strconv.FormatUint(user.Id, 10)
		item["username"] = user.Username
		item["rating"] = strconv.Itoa(int(user.Rating))
		item["image"] = user.Image

		resp = append(resp, item)
	}

	return resp
}

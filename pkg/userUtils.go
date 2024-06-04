package pkg

import (
	"errors"
	"strconv"
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

	review := Review{Reviewer: user, Review: reviewData.Review, Rating: reviewData.Rating}

	if err := DB.Create(&review).Error; err != nil {
		return resp, err
	}

	if err := DB.Model(&toUser).Association("Reviews").Append(&review); err != nil {
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

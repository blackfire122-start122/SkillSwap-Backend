package pkg

import (
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

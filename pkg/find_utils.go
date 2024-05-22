package pkg

import (
	"strconv"
)

func FindUsers(find string) ([]map[string]string, error) {
	var users []User
	usersResp := make([]map[string]string, 0)

	if err := DB.Where("LOWER(username) LIKE ?", "%"+find+"%").Find(&users).Error; err != nil {
		return usersResp, err
	}

	for _, user := range users {
		item := make(map[string]string)

		item["id"] = strconv.FormatUint(user.Id, 10)
		item["username"] = user.Username
		item["rating"] = strconv.Itoa(int(user.Rating))
		item["image"] = user.Image

		usersResp = append(usersResp, item)
	}

	return usersResp, nil
}

func FindCategoriesAll(find string) ([]map[string]string, error) {
	var categories []Category

	categoryResp := make([]map[string]string, 0)

	if err := DB.Where("LOWER(name) LIKE ?", "%"+find+"%").Find(&categories).Error; err != nil {
		return categoryResp, err
	}

	for _, category := range categories {
		item := make(map[string]string)

		item["id"] = strconv.FormatUint(category.Id, 10)
		item["name"] = category.Name

		categoryResp = append(categoryResp, item)
	}

	return categoryResp, nil
}

func FindSkillAll(find string) ([]map[string]string, error) {
	var skills []Skill
	skillResp := make([]map[string]string, 0)

	if err := DB.Where("LOWER(name) LIKE ?", "%"+find+"%").Find(&skills).Error; err != nil {
		return skillResp, err
	}

	for _, skill := range skills {
		item := make(map[string]string)
		item["id"] = strconv.FormatUint(skill.Id, 10)
		item["name"] = skill.Name

		skillResp = append(skillResp, item)
	}

	return skillResp, nil

}

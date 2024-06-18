package pkg

import (
	"strconv"
)

func FindUsers(find string) ([]map[string]string, error) {
	var users []User

	if err := DB.Where("LOWER(username) LIKE ?", "%"+find+"%").Find(&users).Error; err != nil {
		return nil, err
	}

	usersResp := GenerateJsonObjectUsers(users)

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

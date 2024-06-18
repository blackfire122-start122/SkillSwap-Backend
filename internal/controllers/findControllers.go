package controllers

import (
	"SkillSwap/pkg"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FindSkills(c *gin.Context) {
	skillName := strings.ToLower(c.Query("skillName"))

	resp, err := pkg.FindSkillAll(skillName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func FindCategories(c *gin.Context) {
	categoryName := strings.ToLower(c.Query("categoryName"))

	var categories []pkg.Category

	if err := pkg.DB.Where("LOWER(name) LIKE ?", "%"+categoryName+"%").Preload("Skills").Find(&categories).Error; err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]interface{}, 0)

	for _, category := range categories {
		item := make(map[string]interface{})
		item["id"] = strconv.FormatUint(category.Id, 10)
		item["name"] = category.Name

		categorySkills := make([]map[string]string, 0)

		for _, skill := range category.Skills {
			itemSkill := make(map[string]string)

			itemSkill["id"] = strconv.FormatUint(skill.Id, 10)
			itemSkill["name"] = skill.Name

			categorySkills = append(categorySkills, itemSkill)
		}

		item["skills"] = categorySkills

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, resp)
}

func FindAll(c *gin.Context) {
	find := strings.ToLower(c.Query("find"))

	resp := make(map[string]interface{})

	usersResp, err := pkg.FindUsers(find)
	skillResp := make([]map[string]string, 0)
	categoryResp := make([]map[string]string, 0)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(usersResp) <= 5 {
		categoryResp, err = pkg.FindCategoriesAll(find)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if len(usersResp)+len(categoryResp) <= 5 {
		skillResp, err = pkg.FindSkillAll(find)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp["users"] = usersResp
	resp["categories"] = categoryResp
	resp["skills"] = skillResp

	c.JSON(http.StatusOK, resp)
}

func FindUsersOnCategory(c *gin.Context) {
	categoryId := c.Query("categoryId")

	var users []pkg.User

	if err := pkg.DB.Limit(20).Joins("JOIN categories_user ON categories_user.user_id = users.id").
		Where("categories_user.category_id = ?", categoryId).Find(&users).Error; err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	resp := pkg.GenerateJsonObjectUsers(users)

	c.JSON(http.StatusOK, resp)
}

func FindUsersOnSkill(c *gin.Context) {
	skillId := c.Query("skillId")

	var users []pkg.User

	if err := pkg.DB.Limit(20).Joins("JOIN user_skills ON user_skills.user_id = users.id").
		Where("user_skills.skill_id = ?", skillId).Find(&users).Error; err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	resp := pkg.GenerateJsonObjectUsers(users)

	c.JSON(http.StatusOK, resp)
}

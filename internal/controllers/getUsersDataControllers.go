package controllers

import (
	"SkillSwap/pkg"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetImageUser(c *gin.Context) {
	c.File("media/users/" + c.Param("filename"))
}

func GetUser(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := pkg.GetUserData(user)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetUserDataUnauthorized(c *gin.Context) {
	userName := c.Param("userName")

	var user pkg.User
	if err := pkg.DB.Where("username=?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := pkg.GetUserData(user)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetBestPerformers(c *gin.Context) {
	var users []pkg.User

	if err := pkg.DB.Order("rating desc").Limit(20).Find(&users).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := pkg.GenerateJsonObjectUsers(users)

	c.JSON(http.StatusOK, resp)
}

func GetPriceSkills(c *gin.Context) {
	var user pkg.User

	if err := pkg.DB.First(&user, c.Query("userId")).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := pkg.DB.Preload("PricesSkills").Find(&user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]interface{}, 0)

	for _, priceSkill := range user.PricesSkills {
		item := make(map[string]interface{})

		item["Price"] = priceSkill.Price
		item["SkillId"] = priceSkill.SkillId

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, resp)
}

func GetReviews(c *gin.Context) {
	var user pkg.User
	if err := pkg.DB.First(&user, c.Query("userId")).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := pkg.DB.Preload("Reviews").Find(&user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]interface{}, 0)

	for _, review := range user.Reviews {
		if err := pkg.DB.Preload("Reviewer").Find(&review).Error; err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		reviewer := make(map[string]interface{})
		reviewer["id"] = review.Reviewer.Id
		reviewer["username"] = review.Reviewer.Username

		item := make(map[string]interface{})

		item["id"] = review.Id
		item["review"] = review.Review
		item["rating"] = review.Rating
		item["reviewer"] = reviewer

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, resp)
}

//func GetMessages(c *gin.Context) {
//	loginUser, user := CheckSessionUser(c.Request)
//
//	if !loginUser {
//		c.Writer.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	if err := DB.Preload("Messages").Find(&user).Error; err != nil {
//		c.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	resp := make([]map[string]interface{}, 0)
//
//	for _, message := range user.Messages {
//		//if err := DB.Preload("Reviewer").Find(&review).Error; err != nil {
//		//	c.Writer.WriteHeader(http.StatusInternalServerError)
//		//	return
//		//}
//
//		//reviewer := make(map[string]interface{})
//		//reviewer["id"] = review.Reviewer.Id
//		//reviewer["username"] = review.Reviewer.Username
//		//
//		//item := make(map[string]interface{})
//
//		item["id"] = message.ID
//		item["review"] = message.ClientMessage
//
//		resp = append(resp, item)
//	}
//
//	c.JSON(http.StatusOK, resp)
//}

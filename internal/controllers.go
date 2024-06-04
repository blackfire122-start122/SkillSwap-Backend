package internal

import (
	. "SkillSwap/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetImageUser(c *gin.Context) {
	c.File("media/users/" + c.Param("filename"))
}

func GetPriceSkills(c *gin.Context) {
	var user User

	if err := DB.First(&user, c.Query("userId")).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := DB.Preload("PricesSkills").Find(&user).Error; err != nil {
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
	var user User
	if err := DB.First(&user, c.Query("userId")).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := DB.Preload("Reviews").Find(&user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]interface{}, 0)

	for _, review := range user.Reviews {
		if err := DB.Preload("Reviewer").Find(&review).Error; err != nil {
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

func GetUser(c *gin.Context) {
	loginUser, user := CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := GetUserData(user)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetUserDataUnauthorized(c *gin.Context) {
	userName := c.Param("userName")

	var user User
	if err := DB.Where("username=?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := GetUserData(user)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func LoginUser(c *gin.Context) {
	resp := make(map[string]string)

	var user UserLogin
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Password == "" || user.Username == "" {
		resp["Login"] = "Not all field"

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Login(c.Writer, c.Request, &user); err != nil {
		if err.Error() == "record not found" {
			if err := Sign(c.Writer, c.Request, &user); err != nil {
				resp["Login"] = "Error create user"
				c.JSON(http.StatusBadRequest, resp)
				return
			}

			resp["Login"] = "User has been registered successfully"
			c.JSON(http.StatusOK, resp)
		} else {
			resp["Login"] = "Error login"
			c.JSON(http.StatusBadRequest, resp)
		}
	} else {
		resp["Login"] = "OK"
		c.JSON(http.StatusOK, resp)
	}
}

func LogoutUser(c *gin.Context) {
	resp := make(map[string]string)

	if Logout(c.Writer, c.Request) {
		resp["Logout"] = "OK"
		c.JSON(http.StatusOK, resp)
	} else {
		resp["Logout"] = "error logout user"
		c.JSON(http.StatusInternalServerError, resp)
	}
}

func GetBestPerformers(c *gin.Context) {
	var users []User

	if err := DB.Order("rating desc").Limit(10).Find(&users).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]string, 0)

	for _, user := range users {
		item := make(map[string]string)
		item["id"] = strconv.FormatUint(user.Id, 10)
		item["username"] = user.Username
		item["image"] = user.Image
		item["rating"] = strconv.Itoa(int(user.Rating))

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, resp)
}

func ChangeData(c *gin.Context) {
	loginUser, user := CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := make(map[string]string)

	var changedUser ChangedUser
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &changedUser); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ChangeUser(&changedUser, user); err != nil {
		if err.Error() == "user with the same username already exists" {
			resp["Change"] = "User with the same username already exists"
		} else {
			resp["Change"] = "Error change user"
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	} else {
		resp["Change"] = "OK"
		c.JSON(http.StatusOK, resp)
		return
	}
}

func SetUserImage(c *gin.Context) {
	loginUser, user := CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("%d_%s_%s", user.ID, filepath.Base(file.Filename), GenerateSalt())

	for {
		if _, err := os.Stat("media/users/" + fileName); err == nil {
			fileName = fmt.Sprintf("%d_%s_%s", user.ID, filepath.Base(file.Filename), GenerateSalt())
		} else {
			break
		}
	}

	if err := c.SaveUploadedFile(file, "media/users/"+fileName); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Image != "" {
		if err := os.Remove("media/users/" + user.Image); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	user.Image = fileName

	if err := DB.Save(user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func FindSkills(c *gin.Context) {
	skillName := strings.ToLower(c.Query("skillName"))

	var skills []Skill

	if err := DB.Where("LOWER(name) LIKE ?", "%"+skillName+"%").Find(&skills).Error; err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]string, 0)

	for _, skill := range skills {
		item := make(map[string]string)
		item["id"] = strconv.FormatUint(skill.Id, 10)
		item["name"] = skill.Name

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, resp)
}

func FindCategories(c *gin.Context) {
	categoryName := strings.ToLower(c.Query("categoryName"))

	var categories []Category

	if err := DB.Where("LOWER(name) LIKE ?", "%"+categoryName+"%").Preload("Skills").Find(&categories).Error; err != nil {
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

	usersResp, err := FindUsers(find)
	skillResp := make([]map[string]string, 0)
	categoryResp := make([]map[string]string, 0)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(usersResp) <= 5 {
		categoryResp, err = FindCategoriesAll(find)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if len(usersResp)+len(categoryResp) <= 5 {
		skillResp, err = FindSkillAll(find)
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

// ToDo need add limit resp 20

func CreateReview(c *gin.Context) {
	loginUser, user := CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reviewData ReviewData
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &reviewData); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := CreateReviewUser(user, reviewData)

	if err != nil {
		if err.Error() == "rating should be no more than 100 and no less than 0" {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

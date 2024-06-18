package controllers

import (
	"SkillSwap/pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ChangeData(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := make(map[string]string)

	var changedUser pkg.ChangedUser
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &changedUser); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := pkg.ChangeUser(&changedUser, user); err != nil {
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
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("%d_%s_%s", user.ID, filepath.Base(file.Filename), pkg.GenerateSalt())

	for {
		if _, err := os.Stat("media/users/" + fileName); err == nil {
			fileName = fmt.Sprintf("%d_%s_%s", user.ID, filepath.Base(file.Filename), pkg.GenerateSalt())
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

	if err := pkg.DB.Save(user).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func CreateReview(c *gin.Context) {
	loginUser, user := pkg.CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reviewData pkg.ReviewData
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(bodyBytes, &reviewData); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := pkg.CreateReviewUser(user, reviewData)

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

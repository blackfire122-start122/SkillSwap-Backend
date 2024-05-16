package internal

import (
	. "SkillSwap/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func GetImageUser(c *gin.Context) {
	c.File("media/users/" + c.Param("filename"))
}

func GetUser(c *gin.Context) {
	loginUser, user := CheckSessionUser(c.Request)

	if !loginUser {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := make(map[string]string)
	resp["id"] = strconv.FormatUint(user.Id, 10)
	resp["username"] = user.Username
	resp["image"] = user.Image
	resp["email"] = user.Email
	resp["phone"] = user.Phone

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
			if err := Sign(&user); err != nil {
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

	if err := DB.Find(&users).Error; err != nil {
		fmt.Println(err)
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

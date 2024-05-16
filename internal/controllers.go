package internal

import (
	. "SkillSwap/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

	if Login(c.Writer, c.Request, &user) {
		resp["Login"] = "OK"
		c.JSON(http.StatusOK, resp)
	} else {
		if err := Sign(&user); err != nil {
			resp["Login"] = "Error create user"
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		resp["Login"] = "User has been registered successfully"
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

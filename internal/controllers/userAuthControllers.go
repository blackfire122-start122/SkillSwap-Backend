package controllers

import (
	"SkillSwap/pkg"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	resp := make(map[string]string)

	var user pkg.UserLogin
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

	if err := pkg.Login(c.Writer, c.Request, &user); err != nil {
		if err.Error() == "record not found" {
			if err := pkg.Sign(c.Writer, c.Request, &user); err != nil {
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

	if pkg.Logout(c.Writer, c.Request) {
		resp["Logout"] = "OK"
		c.JSON(http.StatusOK, resp)
	} else {
		resp["Logout"] = "error logout user"
		c.JSON(http.StatusInternalServerError, resp)
	}
}

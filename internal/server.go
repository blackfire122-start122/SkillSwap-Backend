package internal

import (
	"github.com/gin-gonic/gin"
)

func SetRouters(router *gin.Engine) {
	router.GET("api/v1/user/image/:filename", GetImageUser) // ToDo nginx
	//router.GET("ws/:roomId", handleConnections)
	router.GET("api/v1/user/getUser", GetUser)
	router.POST("api/v1/user/changeUser", ChangeData)
	router.POST("api/v1/user/setImage", SetUserImage)
	router.POST("api/v1/user/login", LoginUser)
	router.GET("api/v1/user/logout", LogoutUser)
	router.GET("api/v1/bestPerformers", GetBestPerformers)
	router.GET("api/v1/findSkills", FindSkills)
	router.GET("api/v1/findCategories", FindCategories)

}

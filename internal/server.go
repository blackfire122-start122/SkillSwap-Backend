package internal

import (
	"github.com/gin-gonic/gin"
)

func SetRouters(router *gin.Engine) {
	router.GET("api/v1/user/image/:filename", GetImageUser) // ToDo nginx
	router.GET("chat/:roomId", handleConnections)
	router.GET("api/v1/user/getUser", GetUser)
	router.GET("api/v1/user/getUserData/:userName", GetUserDataUnauthorized)
	router.GET("api/v1/user/getPriceSkills", GetPriceSkills)
	router.GET("api/v1/user/getReviews", GetReviews)
	//router.GET("api/v1/user/getMessages", GetMessages)
	router.GET("api/v1/user/customerSkillChats", GetCustomerSkillChats)
	router.GET("api/v1/user/performerSkillChats", GetPerformerSkillChats)
	router.GET("api/v1/user/skillChatMessages", GetSkillChatMessages)
	router.POST("api/v1/user/createReview", CreateReview)
	router.POST("api/v1/user/changeUser", ChangeData)
	router.POST("api/v1/user/setImage", SetUserImage)
	router.POST("api/v1/user/order", Order)
	router.POST("api/v1/user/login", LoginUser)
	router.GET("api/v1/user/logout", LogoutUser)
	router.GET("api/v1/user/findAll", FindAll)
	router.GET("api/v1/bestPerformers", GetBestPerformers)
	router.GET("api/v1/findSkills", FindSkills)
	router.GET("api/v1/findCategories", FindCategories)

}

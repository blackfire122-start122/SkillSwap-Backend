package internal

import (
	"SkillSwap/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetRouters(router *gin.Engine) {
	router.GET("api/v1/user/image/:filename", controllers.GetImageUser) // ToDo nginx
	router.GET("api/v1/user/getUser", controllers.GetUser)
	router.GET("api/v1/user/getUserData/:userName", controllers.GetUserDataUnauthorized)
	router.GET("api/v1/user/getPriceSkills", controllers.GetPriceSkills)
	router.GET("api/v1/user/getReviews", controllers.GetReviews)
	//router.GET("api/v1/user/getMessages", GetMessages)
	router.GET("api/v1/user/customerSkillChats", controllers.GetCustomerSkillChats)
	router.GET("api/v1/user/performerSkillChats", controllers.GetPerformerSkillChats)
	router.GET("api/v1/user/skillChatMessages", controllers.GetSkillChatMessages)
	router.POST("api/v1/user/createReview", controllers.CreateReview)
	router.POST("api/v1/user/changeUser", controllers.ChangeData)
	router.POST("api/v1/user/setImage", controllers.SetUserImage)
	router.POST("api/v1/user/order", controllers.Order)
	router.POST("api/v1/user/login", controllers.LoginUser)
	router.GET("api/v1/user/logout", controllers.LogoutUser)
	router.GET("api/v1/user/findAll", controllers.FindAll)
	router.GET("api/v1/bestPerformers", controllers.GetBestPerformers)
	router.GET("api/v1/findSkills", controllers.FindSkills)
	router.GET("api/v1/findCategories", controllers.FindCategories)
	router.GET("api/v1/findUsersOnCategory", controllers.FindUsersOnCategory)
	router.GET("api/v1/findUsersOnSkill", controllers.FindUsersOnSkill)

	router.GET("chat/:roomId", handleConnections)
}

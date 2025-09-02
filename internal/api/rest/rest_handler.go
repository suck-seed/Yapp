package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/api/rest/handlers"
	"github.com/suck-seed/yapp/internal/services"
)

// RegisterUserRoutes : Group routes proceeding from /user/ and more
func RegisterUserRoutes(r *gin.Engine, userService services.IUserService) {

	// make instance of userHandler
	userHandler := handlers.NewUserHandler(userService)

	// for public
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", userHandler.Ping)
		userGroup.GET("/{user_id}")
		userGroup.GET("/{user_id}/mutual")

		userGroup.POST("/")
	}

	// private (me operation)

	meGroup := r.Group("/me")
	{
		// get my profile
		meGroup.GET("/")
		// update my profile (display, phone, avatar, friend_policy)
		meGroup.PATCH("/")
		// soft delete my profile
		meGroup.DELETE("/")

		meGroup.PATCH("/username")
		meGroup.PATCH("/email")

	}

}

// TODO Make router for halls, messages
func RegisterAuthRoutes(r *gin.Engine, userService services.IUserService) {

	authHandler := handlers.NewAuthHandler(userService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.CreateUser)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/logout", authHandler.Logout)
	}
}

func RegisterHallRoutes(r *gin.Engine, hallService services.IHallService) {

}

func RegisterFloorRoutes(r *gin.Engine, floorService services.IFloorService) {

}

func RegisterRoomRoutes(r *gin.Engine, roomService services.IRoomService) {

}

func RegisterMessageRoutes(r *gin.Engine, messageService services.IMessageService) {

}

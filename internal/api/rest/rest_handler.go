package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/api/rest/handlers"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/services"
)

// TODO Make router for halls, messages
func RegisterAuthRoutes(r *gin.Engine, userService services.IUserService) {

	authHandler := handlers.NewAuthHandler(userService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/signin", authHandler.Signin)
		authGroup.GET("/signout", auth.AuthMiddleware(), authHandler.Signout)
	}
}

// RegisterUserRoutes : Group routes proceeding from /user/ and more
func RegisterUserRoutes(r *gin.RouterGroup, userService services.IUserService) {

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
		meGroup.GET("/", userHandler.GetUser)
		// update my profile (display, phone, avatar, friend_policy)
		meGroup.PATCH("/", userHandler.UpdateUser)
		// soft delete my profile
		meGroup.DELETE("/")
		meGroup.PATCH("/username")
		meGroup.PATCH("/email")
	}

}

func RegisterHallRoutes(r *gin.RouterGroup, hallService services.IHallService) {
	hallHandler := handlers.NewHallHandler(hallService)

	hallGroup := r.Group("/halls")
	{
		hallGroup.GET("/", hallHandler.Ping)
		hallGroup.GET("/{hall_id}")
		hallGroup.GET("/{hall_id}/floors")

		hallGroup.POST("/create", hallHandler.CreateHall)
	}
}

func RegisterFloorRoutes(r *gin.RouterGroup, floorService services.IFloorService) {

}

func RegisterRoomRoutes(r *gin.RouterGroup, roomService services.IRoomService) {

}

func RegisterMessageRoutes(r *gin.RouterGroup, messageService services.IMessageService) {

}

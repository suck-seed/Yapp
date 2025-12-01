package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/api/rest/handlers"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/ws"
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
		userGroup.GET("/:user_id")
		userGroup.GET("/:user_id/mutual")
		userGroup.POST("/")
	}

	// private (me operation)
	meGroup := r.Group("/me")
	{
		// get my profile
		meGroup.GET("/", userHandler.GetUserMe)
		// update my profile (display, phone, avatar, friend_policy)
		meGroup.PATCH("/", userHandler.UpdateUserMe)
		// soft delete my profile
		meGroup.DELETE("/")
		meGroup.PATCH("/username")
		meGroup.PATCH("/email")
	}

}

func RegisterHallRoutes(r *gin.RouterGroup, hallService services.IHallService) {
	hallHandler := handlers.NewHallHandler(hallService)

	halls := r.Group("/halls")
	{

		// TOP LEVEL HALL OPERATIONS
		halls.GET("", hallHandler.GetUserHalls)
		halls.POST("", hallHandler.CreateHall)

		// SINGLE HALL RUD
		halls.GET("/:hallID")
		halls.PATCH("/:hallID")
		halls.DELETE("/:hallID")

		// SETTING SCOPE

		settings := halls.Group("/:hallID/settings")
		{

			// PROFILE MANAGEMENT
			settings.GET("/profile")
			settings.PATCH("/profile")

			// MEMBERS MANAGEMENT
			members := settings.Group("/members")
			{
				members.GET("")
				members.POST("")
				members.PATCH("/:memberID") // updates nickname, timeout, kick, ban, roles, transfer ownership
				members.DELETE("/:memberID")
			}

			// ROLES MANAGEMENT
			roles := settings.Group("/roles")
			{
				roles.GET("")
				roles.POST("")
				roles.PATCH("/:roleID")
				roles.DELETE("/:roleID")

				// roles ko permission
				roles.GET("/:roleID/permissions")
				roles.PATCH("/:roleID/permissions")

			}

			// INVITES MANAGEMENT
			invites := settings.Group("/invites")
			{
				invites.GET("")
				invites.POST("")
				invites.DELETE("/:inviteID") // revoke invite
			}

			// JOIN REQUESTS MANAGEMENT
			requests := settings.Group("/requests")
			{
				requests.GET("")
				requests.POST("")                    // create requests
				requests.PATCH("/:requestID/accept") // accept request
				requests.DELETE("/:requestID")       // accept request
			}

			// BANS
			bans := settings.Group("/bans")
			{
				bans.GET("")
				bans.POST("")          // ban someone
				bans.DELETE("/:banID") // unban
			}
		}
	}
}

func RegisterWebSocketRoutes(r *gin.RouterGroup, hub *ws.Hub, messageService services.IMessageService, hallService services.IHallService, roomService services.IRoomService, userService services.IUserService) {
	wsHandler := ws.NewWebsocketHandler(hub, messageService, hallService, roomService, userService)

	r.GET("/rooms/:room_id", wsHandler.JoinRoom)

}

func RegisterFloorRoutes(r *gin.RouterGroup, floorService services.IFloorService) {

	floorHandler := handlers.NewFloorHandler(floorService)

	floorGroup := r.Group("/floors")
	{
		floorGroup.POST("/", floorHandler.CreateFloor)
	}

}

func RegisterRoomRoutes(r *gin.RouterGroup, roomService services.IRoomService) {

	roomHandler := handlers.NewRoomHandler(roomService)

	roomGroup := r.Group("/rooms")
	{
		roomGroup.POST("/", roomHandler.CreateRoom)
		roomGroup.POST("/add_member")
	}

}

func RegisterMessageRoutes(r *gin.RouterGroup, messageService services.IMessageService) {

	messageHandler := handlers.NewMessageHandler(messageService)

	messageGroup := r.Group("/messages")
	{
		messageGroup.POST("/", messageHandler.FetchMessage)
	}

}

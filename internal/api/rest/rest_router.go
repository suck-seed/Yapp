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
		halls.GET("/:hallID", hallHandler.GetCurrentHall)
		// halls.PATCH("/:hallID", hallHandler.UpdateCurrentHall)
		halls.DELETE("/:hallID", hallHandler.DeleteCurrentHall)

		// SETTING SCOPE

		settings := halls.Group("/:hallID/settings")
		{

			// PROFILE MANAGEMENT
			// RENAME, IMAGE CHANGE, DESCRIPTION CHANGE ETC FROM PROFILE PATCH
			settings.GET("/profile", hallHandler.GetHallProfile)
			settings.PATCH("/profile", hallHandler.UpdateHallProfile)

			// MEMBERS MANAGEMENT
			members := settings.Group("/members")
			{
				members.GET("", hallHandler.GetHallMembers)
				// members.POST("") // There wont be post handler, since we have seperate endpoints for adding and inviting members
				members.PATCH("/:memberID", hallHandler.UpdateHallMember) // updates nickname, timeout, kick, ban, roles, transfer ownership
				members.DELETE("/:memberID", hallHandler.RemoveHallMember)
			}

			// ROLES MANAGEMENT
			roles := settings.Group("/roles")
			{
				roles.GET("", hallHandler.GetHallRoles)
				roles.POST("", hallHandler.CreateHallRoles)
				roles.PATCH("/:roleID", hallHandler.UpdateHallRoles)
				roles.DELETE("/:roleID", hallHandler.DeleteHallRoles)

				// roles ko permission
				roles.GET("/:roleID/permissions", hallHandler.GetRolesPermissions)
				roles.PATCH("/:roleID/permissions", hallHandler.UpdateRolesPermissions)

			}

			// INVITES MANAGEMENT
			invites := settings.Group("/invites")
			{
				invites.GET("", hallHandler.GetCurrentInviteLinks)
				invites.POST("", hallHandler.CreateNewInviteLink)
				invites.DELETE("/:inviteID", hallHandler.InvokeInviteLink) // revoke invite
			}

			// JOIN REQUESTS MANAGEMENT
			requests := settings.Group("/requests")
			{
				requests.GET("", hallHandler.GetCurrentRequests)
				requests.POST("", hallHandler.CreateJoinRequest)                    // create requests
				requests.PATCH("/:requestID/accept", hallHandler.AcceptJoinRequest) // accept request
				requests.DELETE("/:requestID", hallHandler.DeclineJoinRequest)      // accept request
			}

			// BANS
			bans := settings.Group("/bans")
			{
				bans.GET("", hallHandler.GetBannedUsers)
				bans.POST("", hallHandler.BanAnUser)          // ban someone
				bans.DELETE("/:banID", hallHandler.UnbanUser) // unban
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
		floorGroup.POST("", floorHandler.CreateFloor)
	}

}

func RegisterRoomRoutes(r *gin.RouterGroup, roomService services.IRoomService) {

	roomHandler := handlers.NewRoomHandler(roomService)

	roomGroup := r.Group("/rooms")
	{
		roomGroup.POST("", roomHandler.CreateRoom)
		roomGroup.POST("/add_member")
	}

}

func RegisterMessageRoutes(r *gin.RouterGroup, messageService services.IMessageService) {

	messageHandler := handlers.NewMessageHandler(messageService)

	messageGroup := r.Group("/messages")
	{
		messageGroup.POST("", messageHandler.FetchMessage)
	}

}

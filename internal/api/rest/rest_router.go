package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/api/rest/handlers"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/services"
)

// TODO Make router for halls, messages
func RegisterAuthRoutes(r *gin.RouterGroup, userService services.IUserService) {

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

	userGroup := r.Group("/users")
	{
		userGroup.GET("/", userHandler.Ping)
		userGroup.GET("/:user_id", userHandler.GetUserByID)
		userGroup.GET("/:user_id/mutual", auth.AuthMiddleware(), userHandler.GetMutualFriends)
	}

	meGroup := r.Group("/me")
	meGroup.Use(auth.AuthMiddleware())
	{
		meGroup.GET("/", userHandler.GetUserMe)
		meGroup.PATCH("/", userHandler.UpdateUserMe)
		meGroup.DELETE("/", userHandler.DeleteMe)
		meGroup.PATCH("/username", userHandler.UpdateUsername)
		meGroup.PATCH("/email", userHandler.UpdateEmail)

		meGroup.GET("/friends", userHandler.GetMyFriends)
		meGroup.POST("/friends/requests", userHandler.SendFriendRequest)
		meGroup.PATCH("/friends/requests/:request_id", userHandler.RespondFriendRequest)
		meGroup.DELETE("/friends/:user_id", userHandler.Unfriend)

		meGroup.PUT("/app-links", userHandler.UpsertMyAppLink)
		meGroup.DELETE("/app-links/:provider", userHandler.DeleteMyAppLink)
	}

}

func RegisterHallRoutes(r *gin.RouterGroup, hallService services.IHallService, roleServices services.IRoleService, banServices services.IBanService, inviteService services.IInviteService, floorService services.IFloorService, roomService services.IRoomService, messageService services.IMessageService) {
	hallHandler := handlers.NewHallHandler(hallService, roleServices, banServices)
	inviteHandler := handlers.NewInviteHandler(inviteService)

	halls := r.Group("/halls")
	{

		// TOP LEVEL HALL OPERATIONS
		halls.GET("", hallHandler.GetUserHalls)
		halls.POST("", hallHandler.CreateHall)

		// JOIN HALL
		halls.POST("/:hallID/join", hallHandler.JoinHall)

		// SINGLE HALL RUD
		halls.GET("/:hallID", hallHandler.GetCurrentHall)
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
				members.GET("/:memberID", hallHandler.GetHallMember)
				// members.POST("") // There wont be post handler, since we have seperate endpoints for adding and inviting members

				members.PATCH("/:memberID/role", hallHandler.UpdateHallMemberRole)         // updates roles
				members.PATCH("/:memberID/nickname", hallHandler.UpdateHallMemberNickname) // updates nickname
				members.DELETE("/:memberID", hallHandler.KickHallMember)                   // remove member
			}

			// ROLES MANAGEMENT
			roles := settings.Group("/roles")
			{
				roles.GET("", hallHandler.GetHallRoles)
				roles.GET("/:roleID", hallHandler.GetHallRole)
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
				invites.GET("", inviteHandler.ListInviteLinks)
				invites.POST("", inviteHandler.CreateInviteLink)
				invites.DELETE("/:inviteID", inviteHandler.RevokeInviteLink) // revoke invite
			}

			// JOIN REQUESTS MANAGEMENT
			requests := settings.Group("/requests")
			{
				requests.GET("", hallHandler.GetCurrentRequests)
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

		// Halls scoped routes
		hallScoped := halls.Group("/:hallID")
		{
			RegisterFloorRoutes(hallScoped, floorService)
			RegisterRoomRoutes(hallScoped, roomService, messageService)
		}
	}
}

// Separate top-level registration — these routes are NOT under /halls
// because the user only has the code, not a hallID, when clicking the link.
func RegisterInvitePublicRoutes(r *gin.RouterGroup, inviteService services.IInviteService) {
	inviteHandler := handlers.NewInviteHandler(inviteService)

	invites := r.Group("/invites")
	{
		invites.GET("/:code", inviteHandler.GetInviteLinkInfo) // public
	}
}

func RegisterInvitePrivateRoutes(r *gin.RouterGroup, inviteService services.IInviteService) {
	inviteHandler := handlers.NewInviteHandler(inviteService)

	invites := r.Group("/invites")
	{
		invites.POST("/:code/accept", inviteHandler.AcceptInviteLink) // authenticated
	}
}

func RegisterFloorRoutes(r *gin.RouterGroup, floorService services.IFloorService) {

	floorHandler := handlers.NewFloorHandler(floorService)

	floorGroup := r.Group("/floors")
	{
		floorGroup.POST("", floorHandler.CreateFloor)
		floorGroup.GET("", floorHandler.GetFloors) // ?hall_id=
		floorGroup.GET("/:id", floorHandler.GetFloor)
		floorGroup.DELETE("/:id", floorHandler.DeleteFloor)
		floorGroup.PUT("/:id/move", floorHandler.MoveFloor)

		// Private floor member access
		floorGroup.GET("/:id/members", floorHandler.GetFloorMembers)
		floorGroup.GET("/:id/members/:memberID", floorHandler.GetFloorMember)
		floorGroup.PUT("/:id/members/:memberID", floorHandler.AddFloorMember)
		floorGroup.DELETE("/:id/members/:memberID", floorHandler.RemoveFloorMember)
	}

}
func RegisterRoomRoutes(r *gin.RouterGroup, roomService services.IRoomService, messageService services.IMessageService) {
	roomHandler := handlers.NewRoomHandler(roomService)

	roomGroup := r.Group("/rooms")
	{
		roomGroup.POST("", roomHandler.CreateRoom)
		roomGroup.GET("", roomHandler.GetHallRooms)
		roomGroup.GET("/:roomID", roomHandler.GetRoom)
		roomGroup.PATCH("/:roomID", roomHandler.UpdateRoom)
		roomGroup.DELETE("/:roomID", roomHandler.DeleteRoom)
		roomGroup.PUT("/:roomID/move", roomHandler.MoveRoom)

		// Private room member access
		roomGroup.GET("/:roomID/members", roomHandler.GetRoomMembers)
		roomGroup.GET("/:roomID/members/:memberID", roomHandler.GetRoomMember)
		roomGroup.PUT("/:roomID/members/:memberID", roomHandler.AddRoomMember)
		roomGroup.DELETE("/:roomID/members/:memberID", roomHandler.RemoveRoomMember)

		// sync room members to floor members
		roomGroup.PUT("/:roomID/sync-floor-members", roomHandler.SyncRoomMembersToFloor)

		// Room scoped routes
		roomScoped := roomGroup.Group("/:roomID")
		{
			RegisterMessageRoutes(roomScoped, messageService)
		}
	}
}

func RegisterMessageRoutes(r *gin.RouterGroup, messageService services.IMessageService) {

	messageHandler := handlers.NewMessageHandler(messageService)

	messageGroup := r.Group("/messages")
	{
		messageGroup.GET("", messageHandler.FetchMessages)
		messageGroup.GET("/:messageID", messageHandler.GetMessage)
		messageGroup.PATCH("/:messageID", messageHandler.UpdateMessage)
		messageGroup.DELETE("/:messageID", messageHandler.DeleteMessage)
		messageGroup.PUT("/:messageID/reactions/:emoji", messageHandler.AddReaction)
		messageGroup.DELETE("/:messageID/reactions/:emoji", messageHandler.RemoveReaction)
	}

}

func RegisterPresenceRoutes(r *gin.RouterGroup, presenceService services.IPresenceService) {
	presenceHandler := handlers.NewPresenceHandler(presenceService)

	presenceGroup := r.Group("/presence")
	{
		presenceGroup.GET("/me", presenceHandler.GetMyPresence)
		presenceGroup.PATCH("/me", presenceHandler.UpdateMyPresence)
		presenceGroup.GET("/users", presenceHandler.GetManyPresence) // ?ids=id1,id2,id3
	}
}

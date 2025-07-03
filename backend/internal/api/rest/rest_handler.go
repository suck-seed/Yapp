package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/api/rest/handlers"
	"github.com/suck-seed/yapp/internal/services/user"
)

// RegisterUserRoutes : Group routes proceeding from /user/ and more
func RegisterUserRoutes(r *gin.Engine, userService user.IUserService) {

	// make instance of userHandler
	userHandler := handlers.NewUserHandler(userService)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUser)
	}

}

//TODO Make router for halls, messages

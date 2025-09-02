package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Custom error struct

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CustomError) Error() string {
	return e.Message
}

var (
	// AUTH ERRORS
	ErrorUserNotFound   = &CustomError{Code: http.StatusNotFound, Message: "User not found"}
	ErrorEmailExists    = &CustomError{Code: http.StatusConflict, Message: "Email already registered"}
	ErrorUsernameExists = &CustomError{Code: http.StatusConflict, Message: "Username already registered"}
	ErrorNumberExists   = &CustomError{Code: http.StatusConflict, Message: "Number already registered"}
	ErrorWrongPassword  = &CustomError{Code: http.StatusUnauthorized, Message: "Incorrect Password"}

	// REQUEST ERRORS
	ErrorInvalidInput       = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Data"}
	ErrorInvalidUsername    = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Username"}
	ErrorInvalidEmail       = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Email"}
	ErrorInvalidPhoneNumber = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Phone Number"}
	ErrorInvalidDisplayName = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Display Name"}
	ErrorInvalidPassword    = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Password Format"}
	ErrorPasswordWhiteSpace = &CustomError{Code: http.StatusBadRequest, Message: "Password has whitespace"}

	// INTERNAL ERRROS
	ErrorInternal     = &CustomError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrorCreatingUser = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating user"}

	// WEBSOCKET ERRORS
	ErrorFailedUpgrade = &CustomError{Code: http.StatusBadRequest, Message: "Failed to upgrade connection"}
	ErrorInvalidRoomId = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Room Id, can not join room"}

	// CONTEXT ERRORS
	ErrorNoUserIdInContext    = &CustomError{Code: http.StatusBadRequest, Message: "No UserId in context"}
	ErrorEmptyUserIdInContext = &CustomError{Code: http.StatusBadRequest, Message: "Empty UserId in context"}

	// TOKENS
	ErrorMissingToken = &CustomError{Code: http.StatusUnauthorized, Message: "Missing Token"}
	ErrorInvalidToken = &CustomError{Code: http.StatusUnauthorized, Message: "Invalid Token"}
)

// Writting Errors from handlers to client

func WriteError(c *gin.Context, err error) {

	// check if error is defined error
	if customError, ok := err.(*CustomError); ok {

		c.JSON(customError.Code, gin.H{
			"code":    customError.Code,
			"error":   customError.Message,
			"success": false,
		})

		return
	}

	// fallback for unknown errors
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"error":   err.Error(),
		"success": false,
	})

}

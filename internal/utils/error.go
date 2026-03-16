package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	ErrorInvalidUserName    = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Username"}
	ErrorInvalidHallName    = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Hall Name"}
	ErrorInvalidFloorName   = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Floor Name"}
	ErrorInvalidEmail       = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Email"}
	ErrorInvalidPhoneNumber = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Phone Number"}
	ErrorInvalidDisplayName = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Display Name"}
	ErrorInvalidPassword    = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Password Format"}
	ErrorPasswordWhiteSpace = &CustomError{Code: http.StatusBadRequest, Message: "Password has whitespace"}
	ErrorUserDoesntExist    = &CustomError{Code: http.StatusBadRequest, Message: "User Does Not Exist"}

	// ROOM TYPE
	ErrorInvalidRoomType = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Room Type"}

	// INTERNAL ERRROS
	ErrorInternal = &CustomError{Code: http.StatusInternalServerError, Message: "Internal server error"}

	// CREATION ERROR
	ErrorCreatingUser  = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating user"}
	ErrorCreatingHall  = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall"}
	ErrorCreatingFloor = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating Floor"}
	ErrorCreatingRoom  = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating Room"}

	// FETCHING ERROR
	ErrorFetchingRoom       = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Room"}
	ErrorFetchingUser       = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching User"}
	ErrorFetchingMessages   = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Messages"}
	ErrorFetchingHall       = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Hall"}
	ErrorFetchingBan        = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Ban"}
	ErrorFetchingRole       = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Role"}
	ErrorFetchingPermission = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Permissions"}

	// CREATING ERROR
	ErrorCreatingHallRole   = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Role"}
	ErrorCreatingHallMember = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Member"}

	// UPDATING ERROR
	ErrorUpdatingPermissions = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while updating role permissions"}

	// PERMISSIONS DISGARDED ERROR
	ErrorUserCannotManageRoles             = &CustomError{Code: http.StatusUnauthorized, Message: "User does not have privlage to manage roles"}
	ErrorCannotUpdateDefaultRolePermission = &CustomError{Code: http.StatusUnauthorized, Message: "Default Role's Permissions cannot be updated"}
	ErrorCannotUpdateAdminRolePermission   = &CustomError{Code: http.StatusUnauthorized, Message: "Admin Role's Permissions cannot be updated"}

	// ITERATING
	ErrorMessageRowsIteration = &CustomError{Code: http.StatusInternalServerError, Message: "Error occured while iterating message rows"}

	// WEBSOCKET ERRORS
	ErrorFailedUpgrade       = &CustomError{Code: http.StatusBadRequest, Message: "Failed to upgrade connection"}
	ErrorInvalidRoomIDFormat = &CustomError{Code: http.StatusBadRequest, Message: "Invalid Room Id"}

	// CONTEXT ERRORS
	ErrorNoUserIdInContext        = &CustomError{Code: http.StatusBadRequest, Message: "No AuthorID in context"}
	ErrorEmptyUserIdInContext     = &CustomError{Code: http.StatusBadRequest, Message: "Empty AuthorID in context"}
	ErrorInvalidUserIdInContext   = &CustomError{Code: http.StatusUnauthorized, Message: "Bad Token user_id format (not uuid)"}
	ErrorInvalidUsernameInContext = &CustomError{Code: http.StatusUnauthorized, Message: "Bad Token username format"}
	ErrorInvalidUserUUID          = &CustomError{Code: http.StatusBadGateway, Message: "Invalid user UUID in context"}

	// TOKENS
	ErrorMissingToken = &CustomError{Code: http.StatusUnauthorized, Message: "Missing Token"}
	ErrorInvalidToken = &CustomError{Code: http.StatusUnauthorized, Message: "Invalid Token"}
	ErrorTokenExpired = &CustomError{Code: http.StatusUnauthorized, Message: "Token expired"}

	// CURSOR COMBINATION
	ErrorInvalidCursorCombination = &CustomError{Code: http.StatusBadRequest, Message: "Invalid cursor combination, Only 1 cursor is to be sent !"}
	ErrorInvalidCursorLimit       = &CustomError{Code: http.StatusBadRequest, Message: "Invalid cursor Limit, Has to be > 0 !"}

	// DOESNT EXIST
	ErrorRoomNotFound        = &CustomError{Code: http.StatusBadRequest, Message: "Room not found"}
	ErrorHallNotFound        = &CustomError{Code: http.StatusBadRequest, Message: "Hall not found"}
	ErrorRoleNotFound        = &CustomError{Code: http.StatusBadRequest, Message: "Role not found"}
	ErrorPermissionsNotFound = &CustomError{Code: http.StatusBadRequest, Message: "Permissions not found"}
	ErrorBanNotFound         = &CustomError{Code: http.StatusBadRequest, Message: "Ban not found"}
	ErrorFloorNotFound       = &CustomError{Code: http.StatusBadRequest, Message: "Floor not found in this hall"}

	ErrorRequestTimeout = &CustomError{Code: http.StatusRequestTimeout, Message: "Request Timeout"}

	// DOES NOT BELONG
	ErrorUserDoesntBelongRoom       = &CustomError{Code: http.StatusBadRequest, Message: "User is not allowded in this room"}
	ErrorUserDoesntBelongHall       = &CustomError{Code: http.StatusBadRequest, Message: "User is not allowded in this hall"}
	ErrorRoleDoesntBelongInThisHall = &CustomError{Code: http.StatusBadRequest, Message: "Role does not belong to this hall"}
	ErrorInvalidBannerColor         = &CustomError{Code: http.StatusBadRequest, Message: "Invalid banner color"}

	// ALREADY EXIST
	ErrorHallAlreadyExist = &CustomError{Code: http.StatusBadRequest, Message: "Hall under the name already exists"}

	// MESSAGE CREATION ERROR
	ErrorWritingMessage         = &CustomError{Code: http.StatusInternalServerError, Message: "Error Writing Message"}
	ErrorWritingMentions        = &CustomError{Code: http.StatusInternalServerError, Message: "Error Writing Mentions"}
	ErrorFileSizeExceedingLimit = &CustomError{Code: http.StatusUnauthorized, Message: "Uploaded file size exceedes limit"}
	ErrorInvalidFileName        = &CustomError{Code: http.StatusBadRequest, Message: "Invalid file name"}
	ErrorBadFileType            = &CustomError{Code: http.StatusBadRequest, Message: "Bad file type"}
	ErrorFileUnmatch            = &CustomError{Code: http.StatusBadRequest, Message: "File extension does not match in file_type and URL"}
	ErrorLargeFileSize          = &CustomError{Code: http.StatusBadRequest, Message: "File size exceedes allowded size"}
	ErrorCreatingAttachment     = &CustomError{Code: http.StatusInternalServerError, Message: "Error writting attachment to db"}
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

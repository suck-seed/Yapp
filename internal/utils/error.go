package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrorBadContext = &AppError{Code: http.StatusUnauthorized, Message: "Invalid Context"}

	// AUTH ERRORS
	ErrorUserNotFound   = &AppError{Code: http.StatusNotFound, Message: "User not found"}
	ErrorEmailExists    = &AppError{Code: http.StatusConflict, Message: "Email already registered"}
	ErrorUsernameExists = &AppError{Code: http.StatusConflict, Message: "Username already registered"}
	ErrorNumberExists   = &AppError{Code: http.StatusConflict, Message: "Number already registered"}
	ErrorWrongPassword  = &AppError{Code: http.StatusUnauthorized, Message: "Incorrect Password"}

	// REQUEST ERRORS
	ErrorInvalidInput       = &AppError{Code: http.StatusBadRequest, Message: "Invalid Data"}
	ErrorInvalidUserName    = &AppError{Code: http.StatusBadRequest, Message: "Invalid Username"}
	ErrorInvalidHallName    = &AppError{Code: http.StatusBadRequest, Message: "Invalid Hall Name"}
	ErrorInvalidFloorName   = &AppError{Code: http.StatusBadRequest, Message: "Invalid Floor Name"}
	ErrorInvalidEmail       = &AppError{Code: http.StatusBadRequest, Message: "Invalid Email"}
	ErrorInvalidPhoneNumber = &AppError{Code: http.StatusBadRequest, Message: "Invalid Phone Number"}
	ErrorInvalidDisplayName = &AppError{Code: http.StatusBadRequest, Message: "Invalid Display Name"}
	ErrorInvalidPassword    = &AppError{Code: http.StatusBadRequest, Message: "Invalid Password Format"}
	ErrorPasswordWhiteSpace = &AppError{Code: http.StatusBadRequest, Message: "Password has whitespace"}
	ErrorUserDoesntExist    = &AppError{Code: http.StatusBadRequest, Message: "User Does Not Exist"}
	ErrorAlreadyHallMember  = &AppError{Code: http.StatusBadRequest, Message: "User is already this hall's member"}

	// Invalid ID for parsing
	ErrorInvalidIDFormart = &AppError{Code: http.StatusBadRequest, Message: "Error, Invalid ID format"}

	// ROOM TYPE
	ErrorInvalidRoomType = &AppError{Code: http.StatusBadRequest, Message: "Invalid Room Type"}

	// INTERNAL ERRROS
	ErrorInternal = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}

	// CREATION ERROR
	ErrorCreatingUser       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating user"}
	ErrorCreatingHall       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall"}
	ErrorCreatingFloor      = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Floor"}
	ErrorCreatingRoom       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Room"}
	ErrorCreatingHallRole   = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Role"}
	ErrorCreatingHallMember = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Member"}

	// FETCHING ERROR
	ErrorFetchingRoom       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Room Information"}
	ErrorFetchingUser       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching User Information"}
	ErrorFetchingMessages   = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Messages Information"}
	ErrorFetchingHall       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Hall Information"}
	ErrorFetchingFloor      = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Floor Information"}
	ErrorFetchingBan        = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Ban Information"}
	ErrorFetchingRole       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Role Information"}
	ErrorFetchingPermission = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Permissions Information"}

	// UPDATING ERROR
	ErrorUpdatingPermissions = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while updating role permissions"}
	ErrorNoFieldsToUpdate    = &AppError{Code: http.StatusBadRequest, Message: "No field to update"}

	// DELETING ERROR
	ErrorCannotDeleteHall = &AppError{Code: http.StatusUnauthorized, Message: "Only Hall owner can delete hall"}

	// DOESNT EXIST
	ErrorRoomNotFound            = &AppError{Code: http.StatusBadRequest, Message: "Room not found"}
	ErrorHallNotFound            = &AppError{Code: http.StatusBadRequest, Message: "Hall not found"}
	ErrorRoleNotFound            = &AppError{Code: http.StatusBadRequest, Message: "Role not found"}
	ErrorPermissionsNotFound     = &AppError{Code: http.StatusBadRequest, Message: "Permissions not found"}
	ErrorBanNotFound             = &AppError{Code: http.StatusBadRequest, Message: "Ban not found"}
	ErrorFloorNotFound           = &AppError{Code: http.StatusBadRequest, Message: "Floor not found in this hall"}
	ErrorHallDefaultRoleNotFound = &AppError{Code: http.StatusInternalServerError, Message: "Default role not found for hall"}
	ErrorMemberNotFound          = &AppError{Code: http.StatusInternalServerError, Message: "Hall Member not found"}

	// ALREADY EXIST
	ErrorHallAlreadyExist    = &AppError{Code: http.StatusBadRequest, Message: "Hall under the name already exists"}
	ErrorUserAlreadyBanned   = &AppError{Code: http.StatusBadRequest, Message: "User is already banned from this hall"}
	ErrorCannotKickHallOwner = &AppError{Code: http.StatusBadRequest, Message: "Cannot remove the hall owner"}
	ErrorCannotBanHallOwner  = &AppError{Code: http.StatusBadRequest, Message: "Cannot ban the hall owner"}
	ErrorCannotKickYourself  = &AppError{Code: http.StatusBadRequest, Message: "Cannot remove yourself from the hall"}
	ErrorCannotBanYourself   = &AppError{Code: http.StatusBadRequest, Message: "Cannot ban yourself"}

	// DOES NOT BELONG
	ErrorUserDoesntBelongRoom       = &AppError{Code: http.StatusBadRequest, Message: "User is not allowded in this room"}
	ErrorUserDoesntBelongHall       = &AppError{Code: http.StatusBadRequest, Message: "User is not allowded in this hall"}
	ErrorRoleDoesntBelongInThisHall = &AppError{Code: http.StatusBadRequest, Message: "Role does not belong to this hall"}
	ErrorInvalidBannerColor         = &AppError{Code: http.StatusBadRequest, Message: "Invalid banner color"}

	// PERMISSIONS ERROR
	ErrorUserCannotManageRoles             = &AppError{Code: http.StatusUnauthorized, Message: "User does not have privlage to manage roles"}
	ErrorUserCannotKickMembers             = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to kick members"}
	ErrorUserCannotBanMembers              = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to ban members"}
	ErrorUserCannotChangeNickname          = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to change their nickname"}
	ErrorUserCannotManageNicknames         = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to manage nicknames"}
	ErrorCannotUpdateDefaultRolePermission = &AppError{Code: http.StatusUnauthorized, Message: "Default Role's Permissions cannot be updated"}
	ErrorCannotUpdateAdminRolePermission   = &AppError{Code: http.StatusUnauthorized, Message: "Admin Role's Permissions cannot be updated"}
	ErrorCannotModifyProtectedHallRole     = &AppError{Code: http.StatusUnauthorized, Message: "Cannot modify or delete the default or admin hall role"}
	ErrorUserCannotCreateHallRoles         = &AppError{Code: http.StatusUnauthorized, Message: "You are not allowed to create roles in this hall"}
	ErrorRoleNameAlreadyExists             = &AppError{Code: http.StatusConflict, Message: "A role with this name already exists in this hall"}
	ErrorCannotDeleteRoleInUse             = &AppError{Code: http.StatusConflict, Message: "Cannot delete a role that is still assigned to hall members"}
	ErrorUnauthorizedToUpdateHall          = &AppError{Code: http.StatusUnauthorized, Message: "Not Authorized to update hall"}

	// ITERATING
	ErrorMessageRowsIteration = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while iterating message rows"}

	// WEBSOCKET ERRORS
	ErrorFailedUpgrade       = &AppError{Code: http.StatusBadRequest, Message: "Failed to upgrade connection"}
	ErrorInvalidRoomIDFormat = &AppError{Code: http.StatusBadRequest, Message: "Invalid Room Id"}

	// CONTEXT ERRORS
	ErrorNoUserIdInContext        = &AppError{Code: http.StatusBadRequest, Message: "No AuthorID in context"}
	ErrorEmptyUserIdInContext     = &AppError{Code: http.StatusBadRequest, Message: "Empty AuthorID in context"}
	ErrorInvalidUserIdInContext   = &AppError{Code: http.StatusUnauthorized, Message: "Bad Token user_id format (not uuid)"}
	ErrorInvalidUsernameInContext = &AppError{Code: http.StatusUnauthorized, Message: "Bad Token username format"}
	ErrorInvalidUserUUID          = &AppError{Code: http.StatusBadGateway, Message: "Invalid user UUID structure in context"}

	// TOKENS
	ErrorMissingToken = &AppError{Code: http.StatusUnauthorized, Message: "Missing Authorization Token"}
	ErrorInvalidToken = &AppError{Code: http.StatusUnauthorized, Message: "Invalid Authorization Token"}
	ErrorTokenExpired = &AppError{Code: http.StatusUnauthorized, Message: "Authorization Token expired"}

	// CURSOR COMBINATION
	ErrorInvalidCursorCombination = &AppError{Code: http.StatusBadRequest, Message: "Invalid cursor combination, Only 1 cursor is to be sent !"}
	ErrorInvalidCursorLimit       = &AppError{Code: http.StatusBadRequest, Message: "Invalid cursor Limit, Has to be > 0 !"}

	ErrorRequestTimeout = &AppError{Code: http.StatusRequestTimeout, Message: "Request Timeout"}

	// MESSAGE CREATION ERROR
	ErrorWritingMessage         = &AppError{Code: http.StatusInternalServerError, Message: "Error Writing Message"}
	ErrorWritingMentions        = &AppError{Code: http.StatusInternalServerError, Message: "Error Writing Mentions"}
	ErrorFileSizeExceedingLimit = &AppError{Code: http.StatusUnauthorized, Message: "Uploaded file size exceedes limit"}
	ErrorInvalidFileName        = &AppError{Code: http.StatusBadRequest, Message: "Invalid file name"}
	ErrorBadFileType            = &AppError{Code: http.StatusBadRequest, Message: "Bad file type"}
	ErrorFileUnmatch            = &AppError{Code: http.StatusBadRequest, Message: "File extension does not match in file_type and URL"}
	ErrorLargeFileSize          = &AppError{Code: http.StatusBadRequest, Message: "File size exceedes allowded size"}
	ErrorCreatingAttachment     = &AppError{Code: http.StatusInternalServerError, Message: "Error writting attachment to db"}

	// UPDATING ERROR
)

// Writting Errors from handlers to client

func WriteError(c *gin.Context, err error) {

	// check if error is defined error
	if customError, ok := err.(*AppError); ok {

		c.JSON(customError.Code, gin.H{
			"code":    customError.Code,
			"message": customError.Message,
			"success": false,
			"data":    nil,
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

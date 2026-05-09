package utils

import (
	"context"
	"errors"
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
	ErrorTest1 = &AppError{Code: http.StatusForbidden, Message: "Forbidden 1"}
	ErrorTest2 = &AppError{Code: http.StatusForbidden, Message: "Forbidden 2"}
	ErrorTest3 = &AppError{Code: http.StatusForbidden, Message: "Forbidden 3"}
	ErrorTest4 = &AppError{Code: http.StatusForbidden, Message: "Forbidden 4"}

	ErrorForbidden = &AppError{Code: http.StatusForbidden, Message: "Forbidden"}
	// =========================
	// CONTEXT ERRORS
	// =========================
	ErrorBadContext               = &AppError{Code: http.StatusUnauthorized, Message: "Invalid Context"}
	ErrorNoUserIdInContext        = &AppError{Code: http.StatusBadRequest, Message: "No AuthorID in context"}
	ErrorEmptyUserIdInContext     = &AppError{Code: http.StatusBadRequest, Message: "Empty AuthorID in context"}
	ErrorInvalidUserIdInContext   = &AppError{Code: http.StatusUnauthorized, Message: "Bad Token user_id format (not uuid)"}
	ErrorInvalidUsernameInContext = &AppError{Code: http.StatusUnauthorized, Message: "Bad Token username format"}
	ErrorInvalidUserUUID          = &AppError{Code: http.StatusBadGateway, Message: "Invalid user UUID structure in context"}

	// =========================
	// AUTH & TOKEN ERRORS
	// =========================
	ErrorUserNotFound  = &AppError{Code: http.StatusNotFound, Message: "User not found"}
	ErrorWrongPassword = &AppError{Code: http.StatusUnauthorized, Message: "Incorrect Password"}
	ErrorMissingToken  = &AppError{Code: http.StatusUnauthorized, Message: "Missing Authorization Token"}
	ErrorInvalidToken  = &AppError{Code: http.StatusUnauthorized, Message: "Invalid Authorization Token"}
	ErrorTokenExpired  = &AppError{Code: http.StatusUnauthorized, Message: "Authorization Token expired"}

	// =========================
	// CONFLICT / ALREADY EXISTS
	// =========================
	ErrorEmailExists           = &AppError{Code: http.StatusConflict, Message: "Email already registered"}
	ErrorUsernameExists        = &AppError{Code: http.StatusConflict, Message: "Username already registered"}
	ErrorNumberExists          = &AppError{Code: http.StatusConflict, Message: "Number already registered"}
	ErrorHallAlreadyExist      = &AppError{Code: http.StatusBadRequest, Message: "Hall under the name already exists"}
	ErrorRoleNameAlreadyExists = &AppError{Code: http.StatusConflict, Message: "A role with this name already exists in this hall"}
	ErrorUserAlreadyBanned     = &AppError{Code: http.StatusBadRequest, Message: "User is already banned from this hall"}

	// =========================
	// REQUEST VALIDATION ERRORS
	// =========================
	ErrorInvalidInput                    = &AppError{Code: http.StatusBadRequest, Message: "Invalid Data"}
	ErrorCannotBeBothDefaultAndAdminRole = &AppError{Code: http.StatusBadRequest, Message: "A role cannot be both default and admin"}
	ErrorInvalidUserName                 = &AppError{Code: http.StatusBadRequest, Message: "Invalid Username"}
	ErrorInvalidHallName                 = &AppError{Code: http.StatusBadRequest, Message: "Invalid Hall Name"}
	ErrorInvalidFloorName                = &AppError{Code: http.StatusBadRequest, Message: "Invalid Floor Name"}
	ErrorInvalidEmail                    = &AppError{Code: http.StatusBadRequest, Message: "Invalid Email"}
	ErrorInvalidPhoneNumber              = &AppError{Code: http.StatusBadRequest, Message: "Invalid Phone Number"}
	ErrorInvalidDisplayName              = &AppError{Code: http.StatusBadRequest, Message: "Invalid Display Name"}
	ErrorInvalidPassword                 = &AppError{Code: http.StatusBadRequest, Message: "Invalid Password Format"}
	ErrorPasswordWhiteSpace              = &AppError{Code: http.StatusBadRequest, Message: "Password has whitespace"}
	ErrorInvalidIDFormart                = &AppError{Code: http.StatusBadRequest, Message: "Error, Invalid ID format"}
	ErrorInvalidRoomType                 = &AppError{Code: http.StatusBadRequest, Message: "Invalid Room Type"}
	ErrorInvalidBannerColor              = &AppError{Code: http.StatusBadRequest, Message: "Invalid banner color"}
	ErrorInvalidCursorCombination        = &AppError{Code: http.StatusBadRequest, Message: "Invalid cursor combination, Only 1 cursor is to be sent !"}
	ErrorInvalidCursorLimit              = &AppError{Code: http.StatusBadRequest, Message: "Invalid cursor Limit, Has to be > 0 !"}

	// =========================
	// RESOURCE CREATION ERRORS
	// =========================
	ErrorCreatingUser       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating user"}
	ErrorCreatingHall       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall"}
	ErrorCreatingFloor      = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Floor"}
	ErrorCreatingRoom       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Room"}
	ErrorCreatingHallRole   = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Role"}
	ErrorCreatingHallMember = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while creating Hall Member"}

	// =========================
	// RESOURCE FETCHING ERRORS
	// =========================
	ErrorFetchingRoom        = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Room Information"}
	ErrorFetchingUser        = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching User Information"}
	ErrorFetchingMessages    = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Messages Information"}
	ErrorFetchingHall        = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Hall Information"}
	ErrorFetchingHallMembers = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Hall Members"}
	ErrorFetchingFloor       = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Floor Information"}
	ErrorFetchingBan         = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Ban Information"}
	ErrorFetchingRole        = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Role Information"}
	ErrorFetchingPermission  = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while fetching Permissions Information"}

	// =========================
	// MOVING ISSUES
	// =========================
	ErrorMovingRoom = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while Moving Room"}

	// ========================= RESOURCE UPDATING ERRORS
	// =========================
	ErrorUpdatingPermissions = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while updating role permissions"}
	ErrorNoFieldsToUpdate    = &AppError{Code: http.StatusBadRequest, Message: "No field to update"}

	// =========================
	// RESOURCE DELETION ERRORS
	// =========================
	ErrorCannotDeleteHall      = &AppError{Code: http.StatusUnauthorized, Message: "Only Hall owner can delete hall"}
	ErrorCannotDeleteRoleInUse = &AppError{Code: http.StatusUnauthorized, Message: "Cannot delete a role that is still assigned to hall members"}
	ErrorDeletingHall          = &AppError{Code: http.StatusUnauthorized, Message: "Error occured while deleting the hall"}

	// =========================
	// NOT FOUND ERRORS
	// =========================
	ErrorRoomNotFound            = &AppError{Code: http.StatusNotFound, Message: "Room not found"}
	ErrorHallNotFound            = &AppError{Code: http.StatusNotFound, Message: "Hall not found"}
	ErrorRoleNotFound            = &AppError{Code: http.StatusNotFound, Message: "Role not found"}
	ErrorPermissionsNotFound     = &AppError{Code: http.StatusNotFound, Message: "Permissions not found"}
	ErrorBanNotFound             = &AppError{Code: http.StatusNotFound, Message: "Ban not found"}
	ErrorFloorNotFound           = &AppError{Code: http.StatusNotFound, Message: "Floor not found in this hall"}
	ErrorMemberNotFound          = &AppError{Code: http.StatusNotFound, Message: "Hall Member not found"}
	ErrorHallDefaultRoleNotFound = &AppError{Code: http.StatusNotFound, Message: "Hall Member not found"}
	ErrorMessageNotFound         = &AppError{Code: http.StatusNotFound, Message: "Message not found"}
	ErrorReactionNotFound        = &AppError{Code: http.StatusNotFound, Message: "Reaction not found"}

	// =========================
	// JOIN REQUEST ERRORS
	// =========================
	ErrorFetchingJoinRequest = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while fetching Join Request Information"}
	ErrorCreatingJoinRequest = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while creating Join Request"}
	ErrorDeletingJoinRequest = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while deleting Join Request"}
	ErrorJoinRequestNotFound = &AppError{Code: http.StatusNotFound, Message: "Join Request not found"}

	// =========================
	// RELATION / OWNERSHIP ERRORS
	// =========================
	ErrorUserDoesntExist            = &AppError{Code: http.StatusBadRequest, Message: "User Does Not Exist"}
	ErrorAlreadyHallMember          = &AppError{Code: http.StatusBadRequest, Message: "User is already this hall's member"}
	ErrorUserDoesntBelongRoom       = &AppError{Code: http.StatusBadRequest, Message: "User is not allowded in this room"}
	ErrorUserDoesntBelongHall       = &AppError{Code: http.StatusBadRequest, Message: "User is not allowded in this hall"}
	ErrorRoleDoesntBelongInThisHall = &AppError{Code: http.StatusBadRequest, Message: "Role does not belong to this hall"}
	ErrorCannotKickHallOwner        = &AppError{Code: http.StatusBadRequest, Message: "Cannot remove the hall owner"}
	ErrorCannotBanHallOwner         = &AppError{Code: http.StatusBadRequest, Message: "Cannot ban the hall owner"}
	ErrorCannotKickYourself         = &AppError{Code: http.StatusBadRequest, Message: "Cannot remove yourself from the hall"}
	ErrorCannotBanYourself          = &AppError{Code: http.StatusBadRequest, Message: "Cannot ban yourself"}
	ErrorDefaultRoleAlreadyExists   = &AppError{Code: http.StatusBadRequest, Message: "Default role already exists"}
	ErrorAdminRoleAlreadyExists     = &AppError{Code: http.StatusBadRequest, Message: "Admin role already exists"}
	ErrorOwnerRoleAlreadyExists     = &AppError{Code: http.StatusBadRequest, Message: "Owner role already exists"}

	// =========================
	// PERMISSION ERRORS
	// =========================
	ErrorUserCannotManageRoles             = &AppError{Code: http.StatusUnauthorized, Message: "User does not have privlage to manage roles"}
	ErrorUserCannotKickMembers             = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to kick members"}
	ErrorUserCannotBanMembers              = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to ban members"}
	ErrorUserCannotChangeNickname          = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to change their nickname"}
	ErrorUserCannotManageNicknames         = &AppError{Code: http.StatusUnauthorized, Message: "User does not have permission to manage nicknames"}
	ErrorUserCannotCreateHallRoles         = &AppError{Code: http.StatusUnauthorized, Message: "You are not allowed to create roles in this hall"}
	ErrorUserCannotManageInvites           = &AppError{Code: http.StatusUnauthorized, Message: "User does not have privilege to manage invites"}
	ErrorUserCannotManageServer            = &AppError{Code: http.StatusUnauthorized, Message: "User does not have privilege to manage hall"}
	ErrorUserCannotManageRequests          = &AppError{Code: http.StatusUnauthorized, Message: "User does not have privilege to manage requests"}
	ErrorUnauthorizedToUpdateHall          = &AppError{Code: http.StatusUnauthorized, Message: "Not Authorized to update hall"}
	ErrorCannotUpdateDefaultRolePermission = &AppError{Code: http.StatusUnauthorized, Message: "Default Role's Permissions cannot be updated"}
	ErrorCannotUpdateAdminRolePermission   = &AppError{Code: http.StatusUnauthorized, Message: "Admin Role's Permissions cannot be updated"}
	ErrorCannotUpdateHallCreatorsRole      = &AppError{Code: http.StatusUnauthorized, Message: "Cannot update the hall creator's role"}
	ErrorCannotDeleteHallCreatorsRole      = &AppError{Code: http.StatusUnauthorized, Message: "Cannot delete the hall creator's role"}
	ErrorCannotDeleteDefaultRole           = &AppError{Code: http.StatusUnauthorized, Message: "Cannot delete the hall default's role"}
	ErrorNotEnoughPrivlageToDeleteAdmin    = &AppError{Code: http.StatusUnauthorized, Message: "User's role does not have privlate to delete admin role"}
	ErrorCannotModifyAdminRole             = &AppError{Code: http.StatusUnauthorized, Message: "User's role does not have privlage to update admin role"}
	ErrorCannotCreateAdminRole             = &AppError{Code: http.StatusUnauthorized, Message: "User's role does not have privlage to create admin role"}
	ErrorCannotUpdateAdminRole             = &AppError{Code: http.StatusUnauthorized, Message: "User's role does not have privlage to update admin role"}

	// =========================
	// INVITE ERRORS
	// =========================

	// FETCHING
	ErrorFetchingInvite = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while fetching Invite Information"}

	// CREATION / DELETION
	ErrorCreatingInvite       = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while creating Invite"}
	ErrorDeletingInvite       = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while deleting Invite"}
	ErrorGeneratingInviteCode = &AppError{Code: http.StatusInternalServerError, Message: "Error occurred while generating Invite Code"}

	// NOT FOUND
	ErrorInviteNotFound = &AppError{Code: http.StatusNotFound, Message: "Invite not found"}

	// VALIDITY
	ErrorInviteExpired            = &AppError{Code: http.StatusBadRequest, Message: "Invite link has expired"}
	ErrorInviteExhausted          = &AppError{Code: http.StatusBadRequest, Message: "Invite link has reached its maximum uses"}
	ErrorInviteDoesntBelongToHall = &AppError{Code: http.StatusBadRequest, Message: "Invite does not belong to this hall"}

	// REQUEST VALIDATION
	ErrorInvalidExpireAfter = &AppError{Code: http.StatusBadRequest, Message: "Invalid expire_after value"}
	ErrorInvalidMaxUses     = &AppError{Code: http.StatusBadRequest, Message: "Invalid max_uses value"}

	// =========================
	// WEBSOCKET ERRORS
	// =========================
	ErrorFailedUpgrade       = &AppError{Code: http.StatusBadRequest, Message: "Failed to upgrade connection"}
	ErrorInvalidRoomIDFormat = &AppError{Code: http.StatusBadRequest, Message: "Invalid Room Id"}

	// =========================
	// MESSAGE / FILE ERRORS
	// =========================
	ErrorWritingMessage         = &AppError{Code: http.StatusInternalServerError, Message: "Error Writing Message"}
	ErrorWritingMentions        = &AppError{Code: http.StatusInternalServerError, Message: "Error Writing Mentions"}
	ErrorFileSizeExceedingLimit = &AppError{Code: http.StatusUnauthorized, Message: "Uploaded file size exceedes limit"}
	ErrorInvalidFileName        = &AppError{Code: http.StatusBadRequest, Message: "Invalid file name"}
	ErrorBadFileType            = &AppError{Code: http.StatusBadRequest, Message: "Bad file type"}
	ErrorFileUnmatch            = &AppError{Code: http.StatusBadRequest, Message: "File extension does not match in file_type and URL"}
	ErrorLargeFileSize          = &AppError{Code: http.StatusBadRequest, Message: "File size exceedes allowded size"}
	ErrorCreatingAttachment     = &AppError{Code: http.StatusInternalServerError, Message: "Error writting attachment to db"}

	// =========================
	// INTERNAL / SYSTEM ERRORS
	// =========================
	ErrorInternal             = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrorMessageRowsIteration = &AppError{Code: http.StatusInternalServerError, Message: "Error occured while iterating message rows"}
	ErrorRequestTimeout       = &AppError{Code: http.StatusRequestTimeout, Message: "Request Timeout"}

	// conflict / state
	ErrorJoinRequestAlreadyExists      = &AppError{Code: http.StatusConflict, Message: "A pending join request already exists for this hall"}
	ErrorCannotRequestPublicHall       = &AppError{Code: http.StatusBadRequest, Message: "Join requests are only allowed for private halls"}
	ErrorJoinRequestDoesntBelongToHall = &AppError{Code: http.StatusBadRequest, Message: "Join Request does not belong to this hall"}

	// FRIENDSHIP
	ErrorFriendAlreadyExists               = &AppError{Code: http.StatusBadRequest, Message: "Users are already friends"}
	ErrorFriendNotFound                    = &AppError{Code: http.StatusNotFound, Message: "Friend relationship not found"}
	ErrorFriendRequestAlreadyExists        = &AppError{Code: http.StatusBadRequest, Message: "Friend request already exists"}
	ErrorFriendRequestNotFound             = &AppError{Code: http.StatusNotFound, Message: "Friend request not found"}
	ErrorAppLinkNotFound                   = &AppError{Code: http.StatusNotFound, Message: "App link not found"}
	ErrorUnauthorizedToHandleFriendRequest = &AppError{Code: http.StatusNotFound, Message: "Cannot handle other user's Friend Requests"}
)

// Writing Errors from handlers to client

func WriteError(c *gin.Context, err error) {

	if customError, ok := err.(*AppError); ok {
		c.JSON(customError.Code, gin.H{
			"code":    customError.Code,
			"message": customError.Message,
			"success": false,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"error":   err.Error(),
		"success": false,
	})
}

// is Deadline error checker
func IsDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

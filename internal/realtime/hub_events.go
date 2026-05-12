package realtime

import (
	"log"

	"github.com/google/uuid"
)

type HubEventType string

const (
	// User/hall membership events
	HubEventUserJoinedHall     HubEventType = "user_joined_hall"
	HubEventUserLeftHall       HubEventType = "user_left_hall"
	HubEventUserKickedFromHall HubEventType = "user_kicked_from_hall"
	HubEventUserBannedFromHall HubEventType = "user_banned_from_hall"
	HubEventHallDeleted        HubEventType = "hall_deleted"

	// Room events
	HubEventRoomCreated        HubEventType = "room_created"
	HubEventRoomDeleted        HubEventType = "room_deleted"
	HubEventRoomPrivacyChanged HubEventType = "room_privacy_changed"
	HubEventRoomMoved          HubEventType = "room_moved"

	// Direct room access events
	HubEventRoomMemberAdded   HubEventType = "room_member_added"
	HubEventRoomMemberRemoved HubEventType = "room_member_removed"

	// Floor access events
	HubEventFloorMemberAdded    HubEventType = "floor_member_added"
	HubEventFloorMemberRemoved  HubEventType = "floor_member_removed"
	HubEventFloorPrivacyChanged HubEventType = "floor_privacy_changed"
	HubEventFloorDeleted        HubEventType = "floor_deleted"

	// Role/permission events
	HubEventUserAccessResync HubEventType = "user_access_resync"
	HubEventHallAccessResync HubEventType = "hall_access_resync"
)

type HubEvent struct {
	Type HubEventType

	HallID  uuid.UUID
	RoomID  uuid.UUID
	FloorID uuid.UUID

	// UserID is actual users.id, not hall_members.id.
	UserID uuid.UUID

	// MemberID is hall_members.id, useful for service-side context/debug.
	MemberID uuid.UUID

	IsPrivate bool
}

type Publisher interface {
	PublishHubEvent(event HubEvent)
}

type EventBus struct {
	Events chan HubEvent
}

func NewEventBus(buffer int) *EventBus {
	return &EventBus{
		Events: make(chan HubEvent, buffer),
	}
}

func (b *EventBus) PublishHubEvent(event HubEvent) {
	if b == nil {
		return
	}

	select {
	case b.Events <- event:
	default:
		// Do not block REST requests forever if the hub is overloaded.
		// Your fallback is sync_subscriptions/on-demand DB check.
		log.Printf("hub event bus full, dropping event: %+v", event)
	}
}

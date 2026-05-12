package services

import "github.com/suck-seed/yapp/internal/realtime"

func publishHubEvent(p realtime.Publisher, event realtime.HubEvent) {
	if p == nil {
		return
	}
	p.PublishHubEvent(event)
}

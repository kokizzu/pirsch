package model

import (
	"encoding/json"
	"time"
)

// Request represents a request that has or has not been flagged as a bot.
// The creation time, User-Agent, path, and event name are stored in the database to find bots.
type Request struct {
	ClientID    uint64    `db:"client_id" json:"client_id"`
	VisitorID   uint64    `db:"visitor_id" json:"visitor_id"`
	Time        time.Time `json:"time"`
	IP          string    `json:"ip"`
	UserAgent   string    `db:"user_agent" json:"user_agent"`
	Hostname    string    `json:"hostname"`
	Path        string    `json:"path"`
	Event       string    `db:"event_name" json:"event_name"`
	Referrer    string    `json:"referrer"`
	UTMSource   string    `db:"utm_source" json:"utm_source"`
	UTMMedium   string    `db:"utm_medium" json:"utm_medium"`
	UTMCampaign string    `db:"utm_campaign" json:"utm_campaign"`
	Bot         bool      `json:"bot"`
	BotReason   string    `db:"bot_reason" json:"bot_reason"`
}

// String implements the Stringer interface.
func (request Request) String() string {
	out, _ := json.Marshal(request)
	return string(out)
}

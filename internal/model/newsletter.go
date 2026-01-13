package model

// NewsletterFollowRequest represents a request to follow/unfollow a newsletter
type NewsletterFollowRequest struct {
	JID string `json:"jid" example:"123456789@newsletter"`
}

// NewsletterMuteRequest represents a request to mute/unmute a newsletter
type NewsletterMuteRequest struct {
	JID  string `json:"jid" example:"123456789@newsletter"`
	Mute bool   `json:"mute" example:"true"`
}

// NewsletterMarkViewedRequest represents a request to mark messages as viewed
type NewsletterMarkViewedRequest struct {
	JID       string `json:"jid" example:"123456789@newsletter"`
	ServerIDs []int  `json:"server_ids" example:"1,2,3"`
}

// NewsletterReactionRequest represents a request to send a reaction
type NewsletterReactionRequest struct {
	JID      string `json:"jid" example:"123456789@newsletter"`
	ServerID int    `json:"server_id" example:"123"`
	Reaction string `json:"reaction" example:"👍"`
}

// NewsletterCreateRequest represents a request to create a newsletter
type NewsletterCreateRequest struct {
	Name        string `json:"name" example:"My Newsletter"`
	Description string `json:"description" example:"Newsletter description"`
	Picture     string `json:"picture,omitempty"` // base64 encoded image
}

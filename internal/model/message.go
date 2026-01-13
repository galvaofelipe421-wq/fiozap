package model

type TextMessage struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type ImageMessage struct {
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type AudioMessage struct {
	Phone    string `json:"phone"`
	Audio    string `json:"audio"`
	ID       string `json:"id,omitempty"`
	PTT      *bool  `json:"ptt,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type VideoMessage struct {
	Phone    string `json:"phone"`
	Video    string `json:"video"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type DocumentMessage struct {
	Phone    string `json:"phone"`
	Document string `json:"document"`
	FileName string `json:"filename"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type LocationMessage struct {
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
	ID        string  `json:"id,omitempty"`
}

type ContactMessage struct {
	Phone        string `json:"phone"`
	ContactName  string `json:"contact_name"`
	ContactVCard string `json:"contact_vcard"`
	ID           string `json:"id,omitempty"`
}

type ReactionMessage struct {
	Phone     string `json:"phone"`
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

type DeleteMessage struct {
	Phone     string `json:"phone"`
	MessageID string `json:"message_id"`
}

type StickerMessage struct {
	Phone         string   `json:"phone"`
	Sticker       string   `json:"sticker"`
	ID            string   `json:"id,omitempty"`
	PngThumbnail  []byte   `json:"png_thumbnail,omitempty"`
	MimeType      string   `json:"mimetype,omitempty"`
	PackID        string   `json:"pack_id,omitempty"`
	PackName      string   `json:"pack_name,omitempty"`
	PackPublisher string   `json:"pack_publisher,omitempty"`
	Emojis        []string `json:"emojis,omitempty"`
}

type PollMessage struct {
	Phone   string   `json:"phone"`
	Header  string   `json:"header"`
	Options []string `json:"options"`
	ID      string   `json:"id,omitempty"`
}

type ListItem struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	RowID string `json:"row_id,omitempty"`
}

type ListSection struct {
	Title string     `json:"title"`
	Rows  []ListItem `json:"rows"`
}

type ListMessage struct {
	Phone      string        `json:"phone"`
	ButtonText string        `json:"button_text"`
	Title      string        `json:"title"`
	Desc       string        `json:"desc"`
	Sections   []ListSection `json:"sections"`
	FooterText string        `json:"footer_text,omitempty"`
	ID         string        `json:"id,omitempty"`
}

type ButtonItem struct {
	ButtonID   string `json:"button_id"`
	ButtonText string `json:"button_text"`
}

type ButtonsMessage struct {
	Phone   string       `json:"phone"`
	Title   string       `json:"title"`
	Buttons []ButtonItem `json:"buttons"`
	ID      string       `json:"id,omitempty"`
}

type EditMessage struct {
	Phone     string `json:"phone"`
	MessageID string `json:"message_id"`
	Body      string `json:"body"`
}

type MarkReadMessage struct {
	IDs         []string `json:"ids"`
	ChatPhone   string   `json:"chat_phone"`
	SenderPhone string   `json:"sender_phone,omitempty"`
}

type ArchiveChatMessage struct {
	Phone   string `json:"phone"`
	Archive bool   `json:"archive"`
}

type StatusTextMessage struct {
	Body string `json:"body"`
}

type DownloadMediaMessage struct {
	URL           string `json:"url"`
	DirectPath    string `json:"direct_path"`
	MediaKey      []byte `json:"media_key"`
	MimeType      string `json:"mimetype"`
	FileEncSHA256 []byte `json:"file_enc_sha256"`
	FileSHA256    []byte `json:"file_sha256"`
	FileLength    uint64 `json:"file_length"`
}

type GroupPhotoRequest struct {
	GroupJID string `json:"jid"`
	Image    string `json:"image"`
}

type GroupAnnounceRequest struct {
	GroupJID string `json:"jid"`
	Announce bool   `json:"announce"`
}

type GroupLockedRequest struct {
	GroupJID string `json:"jid"`
	Locked   bool   `json:"locked"`
}

type GroupEphemeralRequest struct {
	GroupJID string `json:"jid"`
	Duration string `json:"duration"` // "24h", "7d", "90d", "off"
}

type GroupJoinRequest struct {
	Code string `json:"code"`
}

type GroupInviteInfoRequest struct {
	Code string `json:"code"`
}

type RejectCallRequest struct {
	CallFrom string `json:"call_from"`
	CallID   string `json:"call_id"`
}

type ConnectRequest struct {
	Subscribe []string `json:"subscribe,omitempty"`
	Immediate bool     `json:"immediate,omitempty"`
}

type WebhookRequest struct {
	WebhookURL string   `json:"webhookurl"`
	Events     []string `json:"events,omitempty"`
}

type PairPhoneRequest struct {
	Phone string `json:"phone"`
}

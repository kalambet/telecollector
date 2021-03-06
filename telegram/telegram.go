package telegram

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	EntityTypeBotCommand = "bot_command"
	EntityTypeHashtag    = "hashtag"

	ChatTypeChannel = "channel"

	JoinSeparator = " ➜ "
)

type Chat struct {
	ID               int64            `json:"id"`
	Type             string           `json:"type"`
	Title            string           `json:"title,omitempty"`
	UserName         string           `json:"username,omitempty"`
	FirstName        string           `json:"first_name,omitempty"`
	LastName         string           `json:"last_name,omitempty"`
	Photo            *ChatPhoto       `json:"photo,omitempty"`
	Description      string           `json:"description,omitempty"`
	InviteLink       string           `json:"invite_link,omitempty"`
	PinnedMessage    *Message         `json:"pinned_message,omitempty"`
	Permissions      *ChatPermissions `json:"permissions,omitempty"`
	SlowModeDelay    int              `json:"slow_mode_delay,omitempty"`
	StickerSetName   string           `json:"sticker_set_name,omitempty"`
	CanSetStickerSet bool             `json:"can_set_sticker_set,omitempty"`
}

type User struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	UserName                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

type Response struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url,omitempty"`
	User     *User  `json:"user,omitempty"`
	Language string `json:"language,omitempty"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type ChatPermissions json.RawMessage
type ChatPhoto json.RawMessage
type PhotoSize json.RawMessage
type Audio json.RawMessage
type Video json.RawMessage
type VideoNote json.RawMessage
type Voice json.RawMessage
type Document json.RawMessage
type Animation json.RawMessage
type Game json.RawMessage
type Sticker json.RawMessage
type Contact json.RawMessage
type Venue json.RawMessage
type Poll json.RawMessage
type PollAnswer json.RawMessage
type PassportData json.RawMessage
type SuccessfulPayment json.RawMessage
type Invoice json.RawMessage
type InlineKeyboardMarkup json.RawMessage
type InlineQuery json.RawMessage
type ChosenInlineResult json.RawMessage
type CallbackQuery json.RawMessage
type ShippingQuery json.RawMessage
type PreCheckoutQuery json.RawMessage

type Message struct {
	ID                    int64                 `json:"message_id"`
	From                  *User                 `json:"from,omitempty"`
	Date                  int64                 `json:"date"`
	Chat                  *Chat                 `json:"chat"`
	ForwardFrom           *User                 `json:"forward_from,omitempty"`
	ForwardFromChat       *Chat                 `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID  int64                 `json:"forward_from_message_id,omitempty"`
	ForwardSignature      string                `json:"forward_signature,omitempty"`
	ForwardSenderName     string                `json:"forward_sender_name,omitempty"`
	ForwardDate           int64                 `json:"forward_date,omitempty"`
	ReplyToMessage        *Message              `json:"reply_to_message,omitempty"`
	EditDate              int64                 `json:"edit_date,omitempty"`
	MediaGroupID          string                `json:"media_group_id,omitempty"`
	AuthorSignature       string                `json:"author_signature,omitempty"`
	Text                  string                `json:"text"`
	Entities              []*MessageEntity      `json:"entities"`
	CaptionEntities       []*MessageEntity      `json:"caption_entities"`
	Audio                 *Audio                `json:"audio,omitempty"`
	Document              *Document             `json:"document,omitempty"`
	Animation             *Animation            `json:"animation,omitempty"`
	Game                  *Game                 `json:"game,omitempty"`
	Photo                 []*PhotoSize          `json:"photo_size,omitempty"`
	Sticker               Sticker               `json:"sticker,omitempty"`
	Video                 *Video                `json:"video,omitempty"`
	Voice                 *Voice                `json:"voice,omitempty"`
	VideoNote             *VideoNote            `json:"video_note,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	Contact               *Contact              `json:"contact,omitempty"`
	Location              *Location             `json:"location,omitempty"`
	Venue                 *Venue                `json:"venue,omitempty"`
	Poll                  *Poll                 `json:"poll,omitempty"`
	NewChatMembers        []*User               `json:"new_chat_members,omitempty"`
	LeftChatMember        *User                 `json:"left_chat_member,omitempty"`
	NewChatTitle          string                `json:"new_chat_title,omitempty"`
	NewChatPhoto          []*PhotoSize          `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto       bool                  `json:"delete_chat_photo,omitempty"`
	GroupChatCreated      bool                  `json:"group_chat_created,omitempty"`
	SupergroupChatCreated bool                  `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated    bool                  `json:"channel_chat_created,omitempty"`
	MigrateToChatID       int64                 `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID     int64                 `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage         *Message              `json:"pinned_message,omitempty"`
	Invoice               *Invoice              `json:"invoice,omitempty"`
	SuccessfulPayment     *SuccessfulPayment    `json:"successful_payment,omitempty"`
	ConnectedWebsite      string                `json:"connected_website,omitempty"`
	PassportData          *PassportData         `json:"passport_data,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type Update struct {
	ID                 int64               `json:"update_id"`
	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
}

func (msg *Message) Text2Save() string {
	texts := make([]string, 0)
	if len(msg.Text) != 0 {
		texts = append(texts, msg.Text)
	}

	if msg.ForwardFromChat != nil && msg.ForwardFromChat.Type == ChatTypeChannel {
		texts = append(texts, fmt.Sprintf("https://t.me/%s/%d", msg.ForwardFromChat.UserName, msg.ForwardFromMessageID))
	}

	if len(texts) == 0 {
		return ""
	}

	if len(texts) == 1 {
		return texts[0]
	}

	return strings.Join(texts, JoinSeparator)
}

func (msg *Message) Tags() []string {
	tags := make([]string, 0)
	for _, e := range msg.Entities {
		if e.Type == EntityTypeHashtag {
			tags = append(tags, msg.Text[e.Offset:e.Offset+e.Length])
		}
	}
	return tags
}

func (msg *Message) Command() (string, string) {
	for _, e := range msg.Entities {
		if e.Type == EntityTypeBotCommand {
			// in channels bot command looks like `/command@NameBot`
			// so we split string by @ and then take first segment from second letter to the end
			parts := strings.Split(msg.Text[e.Offset:e.Offset+e.Length], "@")
			var receiver string
			// It could be direct command not in chat
			if len(parts) == 1 {
				receiver = ""
			} else {
				receiver = parts[len(parts)-1]
			}
			return parts[0][1:], receiver
		}
	}

	return "", ""
}

func (msg *Message) Author() *User {
	if msg.From == nil {
		return &User{
			ID:        0,
			FirstName: msg.Chat.Title,
			LastName:  "",
			UserName:  msg.Chat.UserName,
		}
	}

	return msg.From
}

func (user *User) ComposeWhoAmIMessage() string {
	return fmt.Sprintf(
		"*Name*: %s %s\n*Username*: %s\n*ID*:%d",
		user.FirstName,
		user.LastName,
		user.UserName,
		user.ID)
}

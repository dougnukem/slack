package slack

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EventCallback wraps an Event with metadata in the Event API
// see: https://api.slack.com/events-api
type EventCallback struct {
	Token       string          `json:"token,omitempty"`
	TeamID      string          `json:"team_id,omitempty"`
	APIAppID    string          `json:"api_app_id,omitempty"`
	RawEvent    json.RawMessage `json:"event"`
	Event       interface{}     `json:"-"`
	EventTS     string          `json:"event_ts"`
	Type        string          `json:"type"`
	AuthedUsers []string        `json:"authed_users"`
}

type serializedEventCallback EventCallback

// UnmarshalJSON unmarshals an EventCallback
func (e *EventCallback) UnmarshalJSON(b []byte) error {
	// alias type to avoid json unmarshal cycle
	EventCallback := (*serializedEventCallback)(e)
	err := json.Unmarshal(b, EventCallback)
	if err != nil {
		return err
	}

	EventCallback.Event, err = eventMapping.Unmarshal(EventCallback.RawEvent)
	if err != nil {
		return err
	}

	return nil
}

// EventURLVerification is received when you setup the Events API for the first time in your slack App.
// see: https://api.slack.com/events-api#prepare
type EventURLVerification struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

// Event contains the event type
type Event struct {
	Type string `json:"type,omitempty"`
}

// EventUnmarshaller is a mapping of event types to the concrete instances that can be used to unmarshal it
type EventUnmarshaller map[string]interface{}

// Unmarshal will unmarshal an event type
func (e *EventUnmarshaller) Unmarshal(jsonEvent json.RawMessage) (interface{}, error) {
	eventHeader := &Event{}
	err := json.Unmarshal(jsonEvent, eventHeader)
	if err != nil {
		return nil, err
	}

	v, exists := eventMapping[eventHeader.Type]
	if !exists {
		err := fmt.Errorf("Event Error: Received unmapped event %q: %s\n", eventHeader.Type, string(jsonEvent))
		return nil, err
	}

	t := reflect.TypeOf(v)
	event := reflect.New(t).Interface()
	err = json.Unmarshal(jsonEvent, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// eventMapping holds a mapping of event names to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var eventMapping = EventUnmarshaller{
	"message":          MessageEvent{},
	"message.channels": MessageEvent{},
	"message.groups":   MessageEvent{},
	"message.im":       MessageEvent{},
	"message.mpim":     MessageEvent{},

	"presence_change": PresenceChangeEvent{},
	"user_typing":     UserTypingEvent{},

	"channel_marked":          ChannelMarkedEvent{},
	"channel_created":         ChannelCreatedEvent{},
	"channel_joined":          ChannelJoinedEvent{},
	"channel_left":            ChannelLeftEvent{},
	"channel_deleted":         ChannelDeletedEvent{},
	"channel_rename":          ChannelRenameEvent{},
	"channel_archive":         ChannelArchiveEvent{},
	"channel_unarchive":       ChannelUnarchiveEvent{},
	"channel_history_changed": ChannelHistoryChangedEvent{},

	"dnd_updated":      DNDUpdatedEvent{},
	"dnd_updated_user": DNDUpdatedEvent{},

	"im_created":         IMCreatedEvent{},
	"im_open":            IMOpenEvent{},
	"im_close":           IMCloseEvent{},
	"im_marked":          IMMarkedEvent{},
	"im_history_changed": IMHistoryChangedEvent{},

	"group_marked":          GroupMarkedEvent{},
	"group_open":            GroupOpenEvent{},
	"group_joined":          GroupJoinedEvent{},
	"group_left":            GroupLeftEvent{},
	"group_close":           GroupCloseEvent{},
	"group_rename":          GroupRenameEvent{},
	"group_archive":         GroupArchiveEvent{},
	"group_unarchive":       GroupUnarchiveEvent{},
	"group_history_changed": GroupHistoryChangedEvent{},

	"file_created":         FileCreatedEvent{},
	"file_shared":          FileSharedEvent{},
	"file_unshared":        FileUnsharedEvent{},
	"file_public":          FilePublicEvent{},
	"file_private":         FilePrivateEvent{},
	"file_change":          FileChangeEvent{},
	"file_deleted":         FileDeletedEvent{},
	"file_comment_added":   FileCommentAddedEvent{},
	"file_comment_edited":  FileCommentEditedEvent{},
	"file_comment_deleted": FileCommentDeletedEvent{},

	"pin_added":   PinAddedEvent{},
	"pin_removed": PinRemovedEvent{},

	"star_added":   StarAddedEvent{},
	"star_removed": StarRemovedEvent{},

	"reaction_added":   ReactionAddedEvent{},
	"reaction_removed": ReactionRemovedEvent{},

	"pref_change": PrefChangeEvent{},

	"team_join":              TeamJoinEvent{},
	"team_rename":            TeamRenameEvent{},
	"team_pref_change":       TeamPrefChangeEvent{},
	"team_domain_change":     TeamDomainChangeEvent{},
	"team_migration_started": TeamMigrationStartedEvent{},

	"manual_presence_change": ManualPresenceChangeEvent{},

	"user_change": UserChangeEvent{},

	"emoji_changed": EmojiChangedEvent{},

	"commands_changed": CommandsChangedEvent{},

	"email_domain_changed": EmailDomainChangedEvent{},

	"bot_added":   BotAddedEvent{},
	"bot_changed": BotChangedEvent{},

	"accounts_changed": AccountsChangedEvent{},

	"reconnect_url": ReconnectUrlEvent{},

	"event_callback": EventCallback{},

	"url_verification": EventURLVerification{},
}

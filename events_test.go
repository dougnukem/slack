package slack

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newEventCallback(eventJSON string) json.RawMessage {
	return json.RawMessage(
		`
    {
      "token": "z26uFbvR1xHJEdHE1OQiO6t8",
      "team_id": "T061EG9RZ",
      "api_app_id": "A0FFV41KK",
      "event": ` + eventJSON + `,
      "event_ts": "1465244570.336841",
      "type": "event_callback",
      "authed_users": [
              "U061F7AUR"
      ]
    }
    `)
}

const reactionAddedJSON = `{
  "type": "reaction_added",
  "user": "U061F1EUR",
  "item": {
          "type": "message",
          "channel": "C061EG9SL",
          "ts": "1464196127.000002"
  },
  "reaction": "slightly_smiling_face"
}`

var eventCallback = newEventCallback(reactionAddedJSON)

func TestEventCallback(t *testing.T) {
	ev, err := UnmarshalEvent(json.RawMessage(eventCallback))
	assert.NoError(t, err)
	require.IsType(t, &EventCallback{}, ev)

	event := ev.(*EventCallback)
	assert.Equal(t, "event_callback", event.Type)
	assert.Equal(t, "z26uFbvR1xHJEdHE1OQiO6t8", event.Token)
	assert.Equal(t, "T061EG9RZ", event.TeamID)
	assert.Equal(t, "A0FFV41KK", event.APIAppID)
	assert.Equal(t, json.RawMessage(reactionAddedJSON), event.RawEvent)
	assert.Equal(t, "1465244570.336841", event.EventTS)
	require.IsType(t, &ReactionAddedEvent{}, event.Event)
}

const urlVerification = `
{
  "token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
  "challenge": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
  "type": "url_verification"
}
`

func TestEventURLVerification(t *testing.T) {
	ev, err := UnmarshalEvent(json.RawMessage(urlVerification))
	assert.NoError(t, err)
	require.IsType(t, &EventURLVerification{}, ev)

	event := ev.(*EventURLVerification)
	assert.Equal(t, "url_verification", event.Type)
	assert.Equal(t, "Jhj5dZrVaK7ZwHHjRyZWjbDl", event.Token)
	assert.Equal(t, "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P", event.Challenge)
}

func TestEventReactionAdded(t *testing.T) {
	ev, err := UnmarshalEvent(json.RawMessage(eventCallback))
	assert.NoError(t, err)
	require.IsType(t, &EventCallback{}, ev)

	event := ev.(*EventCallback)
	require.IsType(t, &ReactionAddedEvent{}, event.Event)

	innerEvent := event.Event.(*ReactionAddedEvent)
	assert.Equal(t, "reaction_added", innerEvent.Type)
	assert.Equal(t, "U061F1EUR", innerEvent.User)
	assert.Equal(t, "reaction_added", innerEvent.Type)
	assert.Equal(t, "slightly_smiling_face", innerEvent.Reaction)

	assert.Equal(t, "message", innerEvent.Item.Type)
	assert.Equal(t, "C061EG9SL", innerEvent.Item.Channel)
	assert.Equal(t, "1464196127.000002", innerEvent.Item.Timestamp)
}

const messageChannelJSON = `{
    "type": "message.channels",
    "channel": "C2147483705",
    "user": "U2147483697",
    "text": "Hello world",
    "ts": "1355517523.000005"
}`

var messageChannelsEventCallback = newEventCallback(messageChannelJSON)

func TestEventMessageChannels(t *testing.T) {
	ev, err := UnmarshalEvent(json.RawMessage(messageChannelsEventCallback))
	assert.NoError(t, err)
	require.IsType(t, &EventCallback{}, ev)

	event := ev.(*EventCallback)
	require.IsType(t, &MessageEvent{}, event.Event)

	innerEvent := event.Event.(*MessageEvent)
	assert.Equal(t, "message.channels", innerEvent.Type)
	assert.Equal(t, "U2147483697", innerEvent.User)
	assert.Equal(t, "C2147483705", innerEvent.Channel)
	assert.Equal(t, "Hello world", innerEvent.Text)
	assert.Equal(t, "1355517523.000005", innerEvent.Timestamp)
}

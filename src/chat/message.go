package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

var messageBuffer *bytes.Buffer = nil
var messageEncoder *json.Encoder = nil

func init() {
	messageBuffer = bytes.NewBufferString("")
	messageEncoder = json.NewEncoder(messageBuffer)
	messageEncoder.SetEscapeHTML(false)
}

// Message represents a message in the chat
type Message struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

// ToJSON ...
func (t *Message) ToJSON() []byte {
	messageBuffer.Reset()
	if err := messageEncoder.Encode(t); err != nil {
		log.Fatal(err)
	}

	return []byte(messageBuffer.String())
}

// MessageFromJSON unmarshals a Message object from the JSON data
func MessageFromJSON(data []byte) *Message {
	message := &Message{}
	if err := json.Unmarshal(data, message); err != nil {
		log.Fatal(err)
	}

	return message
}

// MessagesToJSON ...
func MessagesToJSON(messages map[string][]*Message) []byte {
	var builder strings.Builder

	builder.WriteString("{")
	chatIdx := 0
	for chatID, chatMsgs := range messages {
		if chatIdx != 0 {
			builder.WriteString(",")
		}

		builder.WriteString("\"")
		builder.WriteString(chatID)
		builder.WriteString("\":")
		builder.WriteString("[")
		for idx, msg := range chatMsgs {
			if idx != 0 {
				builder.WriteString(",")
			}
			builder.WriteString(string(msg.ToJSON()))
		}
		builder.WriteString("]")
		chatIdx++
	}
	builder.WriteString("}")

	return []byte(builder.String())
}

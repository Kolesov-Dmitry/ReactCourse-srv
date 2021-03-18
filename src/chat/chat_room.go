package chat

import (
	"bytes"
	"encoding/json"
	"log"
)

var roomBuffer *bytes.Buffer = nil
var roomEncoder *json.Encoder = nil

func init() {
	roomBuffer = bytes.NewBufferString("")
	roomEncoder = json.NewEncoder(roomBuffer)
	roomEncoder.SetEscapeHTML(false)
}

// Room represents a room in the chat
type Room struct {
	ChatID int64  `json:"chatId"`
	Title  string `json:"title"`
	Income int    `json:"income"`
}

// RoomFromJSON unmarshals a Room object from the JSON data
func RoomFromJSON(data []byte) *Room {
	room := &Room{}
	if err := json.Unmarshal(data, room); err != nil {
		log.Fatal(err)
	}

	return room
}

// ToJSON ...
func (t *Room) ToJSON() []byte {
	roomBuffer.Reset()
	if err := roomEncoder.Encode(t); err != nil {
		log.Fatal(err)
	}

	return []byte(roomBuffer.String())
}

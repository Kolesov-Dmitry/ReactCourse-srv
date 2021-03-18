package chat

import (
	"bytes"
	"encoding/json"
	"log"
)

var profileBuffer *bytes.Buffer = nil
var profileEncoder *json.Encoder = nil

func init() {
	profileBuffer = bytes.NewBufferString("")
	profileEncoder = json.NewEncoder(profileBuffer)
	profileEncoder.SetEscapeHTML(false)
}

// Profile ...
type Profile struct {
	UserName string `json:"userName"`
}

// ToJSON ...
func (t *Profile) ToJSON() []byte {
	profileBuffer.Reset()
	if err := profileEncoder.Encode(t); err != nil {
		log.Fatal(err)
	}

	return []byte(profileBuffer.String())
}

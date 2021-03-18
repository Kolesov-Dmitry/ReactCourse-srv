package net

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"robochat.org/chat-srv/src/chat"
)

// chatsGetHandler '/api/chats' GET request handler
func chatsGetHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	chats := st.LoadRooms()
	chatsJSON, err := json.Marshal(chats)
	if err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write(chatsJSON)

	return nil
}

// messagesGetHandler '/api/messages' GET request handler
func messagesGetHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	msgs := st.LoadMessages()
	msgsJSON := chat.MessagesToJSON(msgs)

	w.Write(msgsJSON)

	return nil
}

type messageType struct {
	ChatID int64  `json:"chatId"`
	MsgID  int64  `json:"msgId"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

// messagesPostHandler '/api/messages' POST request handler
func messagesPostHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	var postMessage messageType
	if err := json.Unmarshal(body, &postMessage); err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	err = st.PushMessage(
		postMessage.ChatID,
		&chat.Message{
			ID:     postMessage.MsgID,
			Author: postMessage.Author,
			Text:   postMessage.Text,
		},
	)

	if err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write([]byte("Success"))

	return nil
}

// messagesPostHandler '/api/messages' DELETE request handler
func messagesDeleteHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	var postMessage messageType
	if err := json.Unmarshal(body, &postMessage); err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	if err := st.DeleteMessage(postMessage.ChatID, postMessage.MsgID); err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write([]byte("Success"))

	return nil
}

type chatRoomType struct {
	ChatID int64  `json:"chatId"`
	Title  string `json:"title"`
}

// chatsPostHandler '/api/chats' POST request handler
func chatsPostHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	var chatRoom chatRoomType
	if err := json.Unmarshal(body, &chatRoom); err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	if err := st.AddRoom(chatRoom.ChatID, chatRoom.Title); err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write([]byte("Success"))

	return nil
}

// chatsDeleteHandler '/api/chats' DELETE request handler
func chatsDeleteHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	var chatRoom chatRoomType
	if err := json.Unmarshal(body, &chatRoom); err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	if err := st.DeleteChat(chatRoom.ChatID); err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write([]byte("Success"))

	return nil
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// profileGetHandler '/api/profile' GET request handler
func profileGetHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	profile := st.LoadProfile()
	w.Write(profile.ToJSON())

	return nil
}

// messagesPostHandler '/api/profile' POST request handler
func profilePostHandler(st chat.Storage, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	var profile chat.Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return &StatusError{err, http.StatusBadRequest}
	}

	if err := st.SetUserName(profile.UserName); err != nil {
		return &StatusError{err, http.StatusInternalServerError}
	}

	w.Write([]byte("Success"))

	return nil
}

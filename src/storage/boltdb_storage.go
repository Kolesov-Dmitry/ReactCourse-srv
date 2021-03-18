package storage

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"robochat.org/chat-srv/src/chat"
)

// BoltDBStorage is catalog.Storage implementation above BoltDB
type BoltDBStorage struct {
	db *bolt.DB
}

// NewBoltDBStorage creates new BoltDBStorage instance
// Inputs:
//  - path to the BoltDB files
func NewBoltDBStorage(p string) (chat.Storage, error) {
	if len(p) == 0 {
		return nil, fmt.Errorf("Database path is not specified")
	}

	db, err := bolt.Open(p, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBStorage{db}, nil
}

// Open ...
func (s *BoltDBStorage) Open() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("Robochat"))
		if err != nil {
			return fmt.Errorf("Could not inititialize database: %v", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte("Chats"))
		if err != nil {
			return fmt.Errorf("Could not inititialize database: %v", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte("Messages"))
		if err != nil {
			return fmt.Errorf("Could not inititialize database: %v", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte("Profile"))
		if err != nil {
			return fmt.Errorf("Could not inititialize database: %v", err)
		}

		return nil
	})
}

// PushMessage ...
func (s *BoltDBStorage) PushMessage(chatID int64, msg *chat.Message) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		chatKey := []byte(strconv.FormatInt(chatID, 10))
		messages := tx.Bucket([]byte("Robochat")).Bucket([]byte("Messages")).Bucket(chatKey)

		messageKey := []byte(strconv.FormatInt(msg.ID, 10))
		if err := messages.Put(messageKey, msg.ToJSON()); err != nil {
			return fmt.Errorf("Could not push message: %v", err)
		}

		return nil
	})
}

// AddRoom ...
func (s *BoltDBStorage) AddRoom(chatID int64, title string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		chats := tx.Bucket([]byte("Robochat")).Bucket([]byte("Chats"))
		messages := tx.Bucket([]byte("Robochat")).Bucket([]byte("Messages"))

		room := &chat.Room{
			ChatID: chatID,
			Title:  title,
			Income: 0,
		}

		chatKey := []byte(strconv.FormatInt(chatID, 10))
		if err := chats.Put(chatKey, room.ToJSON()); err != nil {
			return fmt.Errorf("Could not create chat room: %v", err)
		}

		_, err := messages.CreateBucket(chatKey)
		if err != nil {
			return fmt.Errorf("Could not create chat room: %v", err)
		}

		return nil
	})
}

// LoadRooms ...
func (s *BoltDBStorage) LoadRooms() []*chat.Room {
	result := make([]*chat.Room, 0)

	s.db.View(func(tx *bolt.Tx) error {
		chats := tx.Bucket([]byte("Robochat")).Bucket([]byte("Chats"))

		chats.ForEach(func(key []byte, value []byte) error {
			room := chat.RoomFromJSON(value)
			result = append(result, room)

			return nil
		})

		return nil
	})

	return result
}

// LoadMessages ...
func (s *BoltDBStorage) LoadMessages() map[string][]*chat.Message {
	result := make(map[string][]*chat.Message)

	s.db.View(func(tx *bolt.Tx) error {
		messages := tx.Bucket([]byte("Robochat")).Bucket([]byte("Messages"))

		messages.ForEach(func(key []byte, _ []byte) error {
			messageList := make([]*chat.Message, 0)

			chatBucket := messages.Bucket(key)
			chatBucket.ForEach(func(_ []byte, value []byte) error {
				msg := chat.MessageFromJSON(value)
				messageList = append(messageList, msg)

				return nil
			})

			chatID := string(key)
			result[chatID] = messageList

			return nil
		})

		return nil
	})

	return result
}

// DeleteMessage ...
func (s *BoltDBStorage) DeleteMessage(chatID int64, msgID int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		chatKey := []byte(strconv.FormatInt(chatID, 10))
		messages := tx.Bucket([]byte("Robochat")).Bucket([]byte("Messages")).Bucket(chatKey)

		messageKey := []byte(strconv.FormatInt(msgID, 10))
		if err := messages.Delete(messageKey); err != nil {
			return fmt.Errorf("Could not delete message: %v", err)
		}

		return nil
	})
}

// DeleteChat ...
func (s *BoltDBStorage) DeleteChat(chatID int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		chatKey := []byte(strconv.FormatInt(chatID, 10))
		messages := tx.Bucket([]byte("Robochat")).Bucket([]byte("Messages"))
		chats := tx.Bucket([]byte("Robochat")).Bucket([]byte("Chats"))

		if err := messages.DeleteBucket(chatKey); err != nil {
			return fmt.Errorf("Could not delete chat messages: %v", err)
		}

		if err := chats.Delete(chatKey); err != nil {
			return fmt.Errorf("Could not delete chat: %v", err)
		}

		return nil
	})
}

// LoadProfile ...
func (s *BoltDBStorage) LoadProfile() *chat.Profile {
	userName := ""

	s.db.View(func(tx *bolt.Tx) error {
		profile := tx.Bucket([]byte("Robochat")).Bucket([]byte("Profile"))
		userName = string(profile.Get([]byte("UserName")))

		return nil
	})

	return &chat.Profile{UserName: userName}
}

// SetUserName ...
func (s *BoltDBStorage) SetUserName(userName string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		profile := tx.Bucket([]byte("Robochat")).Bucket([]byte("Profile"))

		if err := profile.Put([]byte("UserName"), []byte(userName)); err != nil {
			return fmt.Errorf("Could not set user name: %v", err)
		}

		return nil
	})
}

// Close ...
func (s *BoltDBStorage) Close() {
	s.db.Close()
}

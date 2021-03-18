package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"time"

	"robochat.org/chat-srv/src/chat"
	"robochat.org/chat-srv/src/net"
	"robochat.org/chat-srv/src/storage"
)

// App represents Catalog Application
type App struct {
	appPath string
	storage chat.Storage
	server  *net.Server
}

// NewApp creates new App instance
func NewApp(args []string) *App {
	return &App{
		appPath: path.Dir(args[0]),
		storage: nil,
		server:  nil,
	}
}

// Init initializes Storage and Server
func (a *App) Init() error {
	var err error

	dbPath := path.Join(a.appPath, "chat.db")
	//dbPath := "d:/Project/ReactCourse-srv/.bin/chat.db"
	if a.storage, err = storage.NewBoltDBStorage(dbPath); err != nil {
		return err
	}

	if err = a.storage.Open(); err != nil {
		return err
	}

	a.server = net.NewServer(5000, a.storage)
	return nil
}

// Run starts the application
func (a *App) Run() {
	// Starting server
	go func() {
		a.server.Start()
	}()

	log.Println("Server started on port: ", a.server.Port())

	a.storage.AddRoom(1, "Room_1")
	a.storage.PushMessage(1, &chat.Message{ID: 1, Author: "User", Text: "Hello Robot"})
	a.storage.PushMessage(1, &chat.Message{ID: 2, Author: "Robot", Text: "Hi User!"})
	a.storage.SetUserName("User")

	// Handle the Interrupt or Kill signals
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill)
	<-exitChan

	log.Println("Shutting down...")
	a.stop()
}

func (a *App) stop() {
	waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Stop(waitCtx); err != nil {
		log.Println(err)
	} else {
		log.Println("Server stopped")
	}

	a.storage.Close()
}

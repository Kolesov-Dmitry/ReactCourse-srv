package main

import (
	"log"
	"os"

	"robochat.org/chat-srv/src/app"
)

func main() {
	a := app.NewApp(os.Args)

	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	a.Run()
}

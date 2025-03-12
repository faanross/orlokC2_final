package main

import (
	"log"
	"orlokC2_final/internal/factory"
)

const serverAddr = ":7777"

func main() {
	listenerFactory := factory.NewListenerFactory()

	l := listenerFactory.CreateListener(serverAddr)

	err := l.Start()
	if err != nil {
		log.Printf("Error starting listener: %s\n", err)
	}
}

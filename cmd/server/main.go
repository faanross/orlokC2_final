package main

import (
	"log"
	"orlokC2_final/internal/factory"
)

var serverAddr = []string{":7777", ":8888", ":9999"}

func main() {
	listenerFactory := factory.NewListenerFactory()

	for _, addr := range serverAddr {

		l := listenerFactory.CreateListener(addr)

		go func(l *factory.Listener) {
			err := l.Start()
			if err != nil {
				log.Printf("Error starting %s: %s\n", l.ID, err)
			}
		}(l)
	}
	select {}
}

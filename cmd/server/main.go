package main

import (
	"fmt"
	"log"
	"orlokC2_final/internal/factory"
	"orlokC2_final/internal/interfaces"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var serverAddr = []string{":7777", ":8888", ":9999"}

func main() {
	stopChan := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	listenerFactory := factory.NewListenerFactory()

	var listeners []*interfaces.Listener

	for _, addr := range serverAddr {
		l := listenerFactory.CreateListener(addr)
		listeners = append(listeners, l)

		go func(l *interfaces.Listener) {
			err := l.Start()

			select {
			case <-stopChan:
				return
			default:
				if err != nil {
					log.Printf("Error starting %s: %s\n", l.ID, err)
				}
			}
		}(l)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Listeners are running. Press Ctrl+C to stop all listeners and exit.")

	<-sigChan

	StopAll(listeners, stopChan)

	fmt.Println("Program exiting gracefully.")
}

func StopAll(listeners []*interfaces.Listener, stopChan chan struct{}) {
	close(stopChan)

	for _, l := range listeners {
		err := l.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", l.ID, err)
		}
	}
	fmt.Println("|STATUS| All listeners shut down.")
}

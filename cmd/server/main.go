package main

import (
	"fmt"
	"log"
	"orlokC2_final/internal/factory"
	"orlokC2_final/internal/listener"
	"orlokC2_final/internal/types"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var listenerConfigs = []struct {
	Addr     string
	Protocol types.ProtocolType
}{
	{":7777", types.H1C}, // HTTP/1.1 on port 7777
	{":8888", types.H1C}, // HTTP/1.1 on port 8888
	{":9999", types.H1C}, // HTTP/1.1 on port 9999
}

func main() {
	stopChan := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	abstractFactory := factory.NewAbstractFactory()

	var listeners []*listener.ConcreteListener

	for _, config := range listenerConfigs {
		l := abstractFactory.CreateListener(config.Protocol, config.Addr)
		listeners = append(listeners, l)

		go func(l *listener.ConcreteListener) {
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

func StopAll(listeners []*listener.ConcreteListener, stopChan chan struct{}) {
	close(stopChan)

	for _, l := range listeners {
		err := l.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", l.ID, err)
		}
	}
	fmt.Println("|STATUS| All listeners shut down.")
}

package run_occasionally

import (
	"time"
	"log"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

const DefaultInterval = "60s"

type OccasionalRunner struct {
	config OccasionalConfiguration
}

func NewOccasionalRunner(config OccasionalConfiguration) *OccasionalRunner {
	runner :=  &OccasionalRunner{
		config: config,
	}

	_, err := time.ParseDuration(config.Interval)
	if err != nil {
		log.Println(fmt.Sprintf("Could not parse interval. Using default: %s", DefaultInterval))
	}

	return runner
}

func (or *OccasionalRunner) Run() {

	commands := make([]*Command, 0)

	duration, err := time.ParseDuration(or.config.Interval)
	if err != nil {
		duration,err = time.ParseDuration(DefaultInterval)
		if err != nil {
			panic(err)
		}
	}

	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			log.Println("Run command")
			command := Command{}
			commands = append(commands, command.Execute(or.config.Command, or.config.Args))
		}
	}()

	var signalWaiter sync.WaitGroup
	signalWaiter.Add(1)

	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	go func() {
		<-signalChannel
		log.Println("Received interruption signal. Exiting...")
		signalWaiter.Done()
	}()

	signalWaiter.Wait()

	for _, command := range commands {
		command.Abort()
	}
}
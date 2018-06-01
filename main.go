package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/robfig/cron"
)

func main() {
	command := os.Getenv("COMMAND")
	interval := os.Getenv("INTERVAL")
	cron := os.Getenv("CRON")

	flag.StringVar(&interval, "interval", interval, "A time.Duration parseable interval. Examples: 5s, 2m or 4h")
	flag.StringVar(&cron, "cron", cron, "A six-field cron expression. Observe the added field for seconds. Example: 54 37 13 * * mon. This will run on every monday at 13:37:54.")
	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		command = strings.Join(args, " ")
	}

	if command == "" && (interval == "" || cron == "") {
		fmt.Printf("usage: %s [-interval i] [-cron expr] command\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	run(command, interval, cron)
}

func run(command string, interval string, cronExpression string) {
	executions := sync.WaitGroup{}

	cron := cron.New()

	var spec string
	if interval != "" {
		spec = fmt.Sprintf("@every %s", interval)
	} else {
		spec = cronExpression
	}

	cron.AddFunc(spec, func() {
		executions.Add(1)
		executeCommand(command)
		executions.Done()
	})

	cron.Start()

	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	<-signalChannel
	log.Println("Waiting for subprocesses to complete execution")
	executions.Wait()
	log.Println("Exiting")
}

func executeCommand(command string) {
	log.Printf("> %s", command)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Unsucessful execution. Statuscode: %d", status.ExitStatus())
				return
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
}

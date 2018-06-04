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
	"github.com/spf13/viper"
)

// JobConfig represents the configuration for a command to be run as
// a periodical job
type Job struct {
	Command  string
	Cron     string
	Interval string
}

// Run the command
func (conf Job) Run() {
	executions.Add(1)
	executeCommand(conf.Command)
	executions.Done()
}

var executions sync.WaitGroup

// main initializes variables and loads configuration
func main() {
	executions = sync.WaitGroup{}

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

	if command != "" && (interval != "" || cron != "") {
		run([]Job{Job{Command: command, Interval: interval, Cron: cron}})
	}

	viper.SetConfigName("run-occasionally")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err == nil {
		var config []Job
		if err = viper.UnmarshalKey("jobs", &config); err == nil {
			run(config)
		}
	} else {
		log.Fatalln(err)
	}

	fmt.Printf("usage: %s [-interval i] [-cron expr] command\n", os.Args[0])
	fmt.Println("Configuration can be loaded from a config file named run-occasionally.[yaml|json] in the current directory.")
	flag.PrintDefaults()
}

// run runs a list of jobs
func run(jobs []Job) {
	daemon := cron.New()

	for _, job := range jobs {
		var spec string
		if job.Interval != "" {
			spec = fmt.Sprintf("@every %s", job.Interval)
		} else {
			spec = job.Cron
		}

		daemon.AddJob(spec, job)
	}

	daemon.Start()

	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	<-signalChannel
	daemon.Stop()
	log.Println("Waiting for subprocesses to complete execution")
	executions.Wait()
	log.Println("Exiting")
	os.Exit(0)
}

// executeCommand executes a command
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

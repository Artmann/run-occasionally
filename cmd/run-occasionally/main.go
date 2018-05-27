package main

import (
	"log"
	"github.com/artmann/run-occasionally"
)

func main() {
	log.Println("Hello World")

	config := run_occasionally.OccasionalConfiguration{
		Command: "/usr/bin/php",
		Args: "/Users/christofferartmann/foo.php",
		Interval: "1s",
	}

	occasionalRunner := run_occasionally.NewOccasionalRunner(config)

	occasionalRunner.Run()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"test/runner"
	"test/watcher"
	"time"

	"github.com/fsnotify/fsnotify"
)

type StringSliceFlag []string

func (s *StringSliceFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *StringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	watchFlag = flag.String("w", "", "Directory to watch")
	commandFlag = flag.String("c", "", "Command to run")
)

func processError(err error) {
	log.Println("error:", err)
}

func main() {
	var ignore StringSliceFlag
	flag.Var(&ignore, "i", "Ignore files")
	flag.Parse()

	dir := *watchFlag
	command := *commandFlag

	if command == "" {
		flag.Usage()
		log.Fatalln("No command supplied")
	}

	if dir == "" {
		dir = "."
		log.Println("No directory supplied - Watching current directory")
	}

	commandParts := strings.Split(command, " ")
	cmd := commandParts[0]
	args := commandParts[1:]

	process := runner.NewRunner(cmd, args, true, func() {
		log.Println("Process stopped")
	}, func(err error) {
		log.Println("Process error: ", err)
	})

	watcher, err := watcher.NewWatcher(dir, true, func(event fsnotify.Event) {
		log.Println("event time:", time.Now().UnixMilli())
		if event.Op&fsnotify.Write == fsnotify.Write {
			fmt.Println("modified file:", event.Name)
			process.Restart()
		}
	}, processError, ignore)

	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Watch(); err != nil {
		log.Fatal(err)
	}

	process.Start()

	registerCleanup(process)
}

func registerCleanup(process runner.Runner) {
	cleanupDone := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan // Block until a signal is received
		fmt.Println("\nReceived an interrupt, performing cleanup...")

		process.Stop()

		cleanupDone <- true
	}()

	fmt.Println("Program is running. Press Ctrl+C to exit...")

	<-cleanupDone // Wait for the cleanup to complete
	fmt.Println("Cleanup completed. Exiting program.")
}
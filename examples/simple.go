package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MorpheusXAUT/zkillredisqo"
	"log"
	"os"
)

var (
	logFileKills  *string
	logKills      bool = false
	logFileErrors *string
	logErrors     bool = false
)

func init() {
	logFileKills = flag.String("kills", "", "File to log retrieved kills into")
	logFileErrors = flag.String("errors", "", "File to log errors into")
	flag.Parse()

	if len(*logFileKills) > 0 {
		logKills = true
	}
	if len(*logFileErrors) > 0 {
		logErrors = true
	}
}

func main() {
	poller := zkillredisqo.NewPoller(nil)
	poller.SetTimeToWait(5)

	for {
		select {
		case kill, ok := <-poller.Kills:
			if !ok {
				log.Fatal("Failed to read from kills channel")
			}

			log.Printf("%+v\n", kill)
			if logKills {
				logKillToFile(kill)
			}
			break
		case err, ok := <-poller.Errors:
			if !ok {
				log.Fatal("Failed to read from errors channel")
			}

			log.Printf("*** ERROR: %v\n", err)
			if logErrors {
				logErrorToFile(err)
			}
			break
		}
	}
}

func logKillToFile(kill *zkillredisqo.Kill) {
	f, err := os.OpenFile(*logFileKills, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("*** FAILED OPEN KILL LOG: %v\n", err)
		return
	}

	defer f.Close()

	data, err := json.Marshal(kill)
	if err != nil {
		log.Printf("*** FAILED MARSHAL KILL: %v\n", err)
		return
	}

	if _, err = f.WriteString(fmt.Sprintf("%s\n", data)); err != nil {
		log.Printf("*** FAILED WRITE KILL LOG: %v\n", err)
		return
	}

	return
}

func logErrorToFile(e error) {
	f, err := os.OpenFile(*logFileErrors, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("*** FAILED OPEN ERROR LOG: %v\n", err)
		return
	}

	defer f.Close()
	if _, err = f.WriteString(fmt.Sprintf("%v\n", e)); err != nil {
		log.Printf("*** FAILED WRITE ERROR LOG: %v\n", err)
		return
	}

	return
}

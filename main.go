package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/theverything/reminder/pkg/reminder"
)

func main() {
	configPath := flag.String("config", "", "path to reminder config")

	flag.Parse()

	if len(*configPath) == 0 {
		panic("missing config path")
	}

	f, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var c reminder.Config
	err = json.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}

	r := reminder.New(c)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		r.Stop()
	}()

	printConfig, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(`Starting reminder with this config`)
	fmt.Println(string(printConfig))

	r.Start()
}

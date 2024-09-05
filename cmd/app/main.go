package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var name, environment string

func init() {
	os.Setenv("TZ", "Africa/Nairobi")
	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	flag.StringVar(&name, "name", "", "Name to greet")
	flag.StringVar(&environment, "environment", "", "Environment")
	flag.Parse()
}

func main() {
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		fmt.Println(err)
	}
	// handle err
	time.Local = loc // -> this is setting the global timezone

	log.Info("time is ", time.Now())

	if name == "" {
		fmt.Println("Please provide a name using the -name flag")
		log.Info("Name is ", name)
		return
	}

	if name != "" {

	} else {
		log.Warn("no environment set!!")
	}
	fmt.Printf("Hello, %s!\n", name)
}

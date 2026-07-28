package main

import (
	"flag"
	"log"
	"time"

	"github.com/babilon15/trfeed/internal/scan"
	"github.com/babilon15/trfeed/pkg/utils"
)

const (
	counts = 3
)

func main() {
	useJSON := flag.Bool("json", false, "The configuration file will be in JSON format instead of YAML.")
	flag.Parse()

	if utils.IsSuperuserNow() {
		log.Fatalln("do not run this program with superuser privileges")
	}

	s := &scan.Scanner{}
	s.Init(*useJSON)

	counter := 0

	for {
		s.Run()
		s.AddHits()
		s.Save()
		time.Sleep(time.Minute)
		if s.Conf.ReloadConfigFile {
			if counter == counts {
				s.GetConfigFile()
				counter = 0
			} else {
				counter++
			}
		}
	}
}

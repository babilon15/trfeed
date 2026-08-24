package main

import (
	"flag"
	"log"
	"time"

	"github.com/babilon15/trfeed/internal/scan"
	"github.com/babilon15/trfeed/pkg/utils"
)

func main() {
	useWebr := flag.Bool("web", false, "Enable web-based status reports.")
	flag.Parse()

	if utils.IsSuperuserNow() {
		log.Fatalln("do not run this program with superuser privileges")
	}

	s := &scan.Scanner{}
	s.Init()

	for {
		s.GetConfigFile()
		s.Run()
		s.AddHits()
		s.Save()
		time.Sleep(time.Minute)
	}
}

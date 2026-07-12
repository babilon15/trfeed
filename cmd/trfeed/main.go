package main

import (
	"log"
	"time"

	"github.com/babilon15/trfeed/internal/scan"
	"github.com/babilon15/trfeed/pkg/utils"
)

func main() {
	if utils.IsSuperuserNow() {
		log.Fatalln("do not run this program with superuser privileges")
	}

	s := &scan.Scanner{}
	s.Init()

	counter := 0

	for {
		s.Run()
		s.AddHits()
		s.Save()
		time.Sleep(time.Minute)
		if s.Conf.ReloadConfigFile {
			if counter == 5 {
				s.GetConfigFile()
				counter = 0
			} else {
				counter++
			}
		}
	}
}

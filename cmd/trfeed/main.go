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

	for {
		s.Run()
		s.Save()
		s.AddHits()
		time.Sleep(time.Minute)
	}
}

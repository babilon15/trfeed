package main

import (
	"fmt"
	"os"
	"time"

	"github.com/babilon15/trfeed/internal/scan"
	"github.com/babilon15/trfeed/pkg/utils"
)

func main() {
	if utils.IsSuperuserNow() {
		fmt.Println("do not run this program with superuser privileges")
		os.Exit(1)
	}

	s := &scan.Scanner{}
	s.Init()

	for {
		s.Run()
		if s.Conf.ConfigOverwrite {
			s.Save()
		}
		s.AddHits()
		time.Sleep(time.Second * 90)
	}
}

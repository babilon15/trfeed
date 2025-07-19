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
		fmt.Println("Do not run this program with superuser privileges.")
		os.Exit(1)
	}

	for {
		scan.Scan()
		scan.AddHits()
		time.Sleep(time.Second * 30)
	}
}

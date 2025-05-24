package main

import (
	"log"

	"github.com/babilon15/trfeed/internal/places"
)

func main() {
	log.SetFlags(0)
	places.Checks()
}

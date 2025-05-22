package main

import (
	"fmt"
	"path"

	"github.com/babilon15/trfeed/internal/places"
)

func main() {
	places.Set()

	fmt.Println(places.ProgramDirParent)
	fmt.Println(places.ConfigFile)

	fmt.Println(path.Join(places.ProgramDirParent, places.ConfigFile))

	fmt.Println(places.ProgramTempDir)

	fmt.Println(places.GetAbs(places.RemnantsFile))

	places.Checks()
}

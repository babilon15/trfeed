package places

import (
	"log"
	"os"
	"os/user"
	"path"
)

const (
	programDirName   = "trfeed"
	torrentsDirName  = "torrents"
	configFileName   = "config.yml"
	remnantsFileName = "remnants.yml"
)

var (
	programDir   string = ""
	ConfigFile   string = ""
	RemnantsFile string = ""
	TorrentsDir  string = ""
)

// Call it before using the variables!
func Checks() {
	// Get the current user:
	currentUser, currentUserErr := user.Current()
	if currentUserErr != nil {
		log.Fatalf(currentUserErr.Error())
	}

	if currentUser.Uid == "0" {
		log.Fatalf("Oh, don't be silly. It's completely unnecessary to do this. ;)")
	}

	// Update those variables:
	programDir = path.Join(currentUser.HomeDir, ".config", programDirName)
	ConfigFile = path.Join(programDir, configFileName)
	RemnantsFile = path.Join(programDir, remnantsFileName)
	TorrentsDir = path.Join(programDir, torrentsDirName)

	// CHECKS:

	if err := os.MkdirAll(TorrentsDir, 0o700); err != nil {
		log.Fatalf(err.Error())
	}
}

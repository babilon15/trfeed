package places

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/babilon15/trfeed/pkg/utils"
	"golang.org/x/sys/unix"
)

const (
	programDirEnvVarName     = "TRFEED_PROG_DIR"
	programTempDirEnvVarName = "TRFEED_TEMP_DIR"
)

var (
	programDirName   string = "trfeed"
	ProgramDirParent string = ""
	ConfigFile       string = path.Join(programDirName, "config.yaml")
	RemnantsFile     string = path.Join(programDirName, "remnants.yaml")
	ProgramTempDir   string = path.Join("/", "tmp", programDirName)
)

func Set() {
	cUsr, cUsrErr := user.Current()
	if cUsrErr != nil {
		log.Fatalf(cUsrErr.Error())
	}

	programDirEnvVar := os.Getenv(programDirEnvVarName)
	if programDirEnvVar != "" {
		ProgramDirParent = programDirEnvVar
	} else {
		ProgramDirParent = path.Join(cUsr.HomeDir, ".config")
	}

	programTempDirEnvVar := os.Getenv(programTempDirEnvVarName)
	if programTempDirEnvVar != "" {
		ProgramTempDir = programTempDirEnvVar
	}
}

func GetAbs(p string) string {
	return path.Join(ProgramDirParent, p)
}

func Checks() {
	if err := unix.Access(ProgramDirParent, unix.W_OK); err != nil {
		log.Fatalf(err.Error())
	}

	if err := os.MkdirAll(GetAbs(programDirName), utils.DMode); err != nil {
		log.Fatalf(err.Error())
	}
}

package utils

import (
	"os/user"
)

func IsSuperuserNow() bool {
	u, _ := user.Current()
	return u.Uid == "0"
}

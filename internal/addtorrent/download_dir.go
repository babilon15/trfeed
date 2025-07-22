package addtorrent

import (
	"os/exec"
	"strings"
)

const (
	ddLine = "Download directory"
)

func GetDownDirWithRemote(host, auth string) string {
	args := []string{}

	if len(host) != 0 {
		args = append(args, host)
	}

	if len(auth) != 0 {
		args = append(args, "--auth", auth)
	}

	args = append(args, "--session-info")

	cmd := exec.Command("transmission-remote", args...)

	out, _ := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		pair := strings.Split(line, ":")
		if len(pair) != 2 {
			continue
		}

		if strings.TrimSpace(pair[0]) == ddLine {
			return strings.TrimSpace(pair[1])
		}
	}

	return ""
}

package addtorrent

import (
	"os/exec"
	"strings"
)

func AddTorrentWithRemote(host, auth, resource, targetDir string, labels []string, paused bool) error {
	args := []string{host}

	if len(auth) >= 3 {
		args = append(args, "--auth", auth)
	}

	args = append(args, "--add", resource)

	if targetDir != "" {
		args = append(args, "--download-dir", targetDir)
	}

	if len(labels) != 0 {
		args = append(args, "--labels", strings.Join(labels, ","))
	}

	if paused {
		args = append(args, "--start-paused")
	} else {
		args = append(args, "--no-start-paused")
	}

	cmd := exec.Command("transmission-remote", args...)

	_, err := cmd.CombinedOutput()

	return err
}

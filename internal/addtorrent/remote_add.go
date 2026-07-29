package addtorrent

import (
	"os/exec"
	"strings"
)

func AddTorrentWithRemote(binPath, host, auth, resource, targetDir string, labels []string, pause bool) error {
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

	if pause {
		args = append(args, "--stop")
	} else {
		args = append(args, "--start")
	}

	if binPath == "" {
		binPath = "transmission-remote"
	}

	cmd := exec.Command(binPath, args...)

	_, err := cmd.CombinedOutput()

	return err
}

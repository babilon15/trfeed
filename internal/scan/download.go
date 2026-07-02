package scan

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url, target, mediaType string) (string, error) {
	resp, respErr := http.Get(url)
	if respErr != nil {
		return "", respErr
	}
	defer resp.Body.Close()

	m := resp.Header.Get("content-type")
	if m != mediaType {
		return "", fmt.Errorf("unexpected Content-Type: expected %q, got %q", mediaType, m)
	}

	_, params, parseMediaErr := mime.ParseMediaType(resp.Header.Get("content-disposition"))
	if parseMediaErr != nil {
		return "", parseMediaErr
	}

	path := filepath.Join(target, params["filename"])

	file, fileErr := os.Create(path)
	if fileErr != nil {
		return "", fileErr
	}
	defer file.Close()

	_, err := io.Copy(file, resp.Body)
	return path, err
}

func DownloadTorrentFile(url, target string) (string, error) {
	return DownloadFile(url, target, "application/x-bittorrent")
}

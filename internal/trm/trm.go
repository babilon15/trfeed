package trm

import (
	"context"
	"net/url"

	"github.com/hekmon/transmissionrpc/v3"
)

// Hozzáadandó torrent a végleges adatokkal:
type TrAdd struct {
	Paused       bool   // Ne induljon el a letöltés a hozzáadás után.
	Path         string // A letöltött torrent fájl elérési útja.
	URL          string // Tartalékként a letöltés URL-je, arra az esetre, ha a letöltött fájl (Path) hozzáadása sikertelen lenne.
	DownloadDir  string // Nincs további ellenőrzés, ha nem sikerült volna az útvonal létrehozása a Transmission még megoldhatja, ha rendelkezik a szükséges jogosultságokkal.
	MainCat      string // -> Labels
	FriendlyName string // Csak ha van. (Labels)
}

func TrConn(endpoint string) (*transmissionrpc.Client, error) {
	p, pErr := url.Parse(endpoint)
	if pErr != nil {
		return nil, pErr
	}
	return transmissionrpc.New(p, nil)
}

func AddTorrent(c *transmissionrpc.Client, t TrAdd) error {
	labels := []string{"ncad"}
	if t.MainCat != "" {
		labels = append(labels, t.MainCat)
	}
	if t.FriendlyName != "" {
		labels = append(labels, t.FriendlyName)
	}
	payload := transmissionrpc.TorrentAddPayload{
		DownloadDir: &t.DownloadDir,
		Filename:    &t.Path,
		Labels:      labels,
		Paused:      &t.Paused,
	}
	_, err := c.TorrentAdd(context.TODO(), payload)
	if err != nil {
		payload.Filename = &t.URL
		_, err = c.TorrentAdd(context.TODO(), payload)
		return err
	}
	return nil
}

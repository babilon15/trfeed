package feed

import (
	"encoding/xml"
	"hash/fnv"
	"log"
	"time"
)

const (
	customLayout = "Mon, 2 Jan 2006 15:04:05 -0700"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Item    []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Category    string   `xml:"category"`
	PubDate     string   `xml:"pubDate"`
	UniqueNum   uint64
}

func (i *Item) ParsePubDate() (time.Time, error) { return time.Parse(customLayout, i.PubDate) }

func (i *Item) GetPubDate() string {
	pubDate, err := i.ParsePubDate()
	if err != nil {
		log.Println(err)
	}

	return pubDate.Format(time.DateTime)
}

func (i *Item) GetUniqueNum() uint64 {
	if i.UniqueNum != 0 {
		return i.UniqueNum
	}

	n64a := fnv.New64a()
	n64a.Write([]byte(i.Title + i.PubDate))
	i.UniqueNum = n64a.Sum64()

	return i.UniqueNum
}

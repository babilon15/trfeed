package feed

import (
	"encoding/xml"
	"hash/fnv"
	"time"
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
	UniqueNum   uint32
}

func (i *Item) ParsePubDate() (time.Time, error) { return time.Parse(time.RFC1123Z, i.PubDate) }

func (i *Item) GetUniqueNum() uint32 {
	if i.UniqueNum != 0 {
		return i.UniqueNum
	}

	n32 := fnv.New32()
	n32.Write([]byte(i.Title + i.PubDate))
	i.UniqueNum = n32.Sum32()

	return i.UniqueNum
}

package in

import (
	"time"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

type FeedFormat interface {
	Namespaces() feeds.Namespaces
	ID() string
	Title() string
	Alternate() string
	Self() string
	Authors() []*feeds.Person
	Contributors() []*feeds.Person
	Categories() []*feeds.Category
	Copyright() string
	Description() string
	Image() string
	Language() string
	Links() []*feeds.Link
	Published() *time.Time
	Updated() *time.Time
	Entries() []*feeds.Entry
}

func elementText(xml *etree.Element, tag string) string {
	return xml.SelectElement(tag).NotNil().Text()
}

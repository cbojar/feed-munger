package in

import (
	"time"

	"github.com/cbojar/feed-munger/feeds"
)

type Feed struct {
	ID           string
	Namespaces   feeds.Namespaces
	Title        string
	Alternate    string
	Self         string
	Authors      []*feeds.Person
	Contributors []*feeds.Person
	Categories   []*feeds.Category
	Copyright    string
	Description  string
	Image        string
	Language     string
	Links        []*feeds.Link
	Published    *time.Time
	Updated      *time.Time
	Entries      []*feeds.Entry
}

func NewFeed(format FeedFormat) *Feed {
	return &Feed{
		ID:           format.ID(),
		Namespaces:   format.Namespaces(),
		Title:        format.Title(),
		Alternate:    format.Alternate(),
		Self:         format.Self(),
		Authors:      format.Authors(),
		Contributors: format.Contributors(),
		Categories:   format.Categories(),
		Copyright:    format.Copyright(),
		Description:  format.Description(),
		Image:        format.Image(),
		Language:     format.Language(),
		Links:        format.Links(),
		Published:    format.Published(),
		Updated:      format.Updated(),
		Entries:      format.Entries(),
	}
}

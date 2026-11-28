package feeds

import (
	"time"

	"github.com/beevik/etree"
)

type Entry struct {
	ID           string
	Title        string
	Updated      *time.Time
	Published    *time.Time
	Links        []*Link
	Summary      string
	Content      string
	Rights       string
	Authors      []*Person
	Contributors []*Person
	Source       string
	Categories   []*Category
	Extensions   []*etree.Element
}

func (entry *Entry) Timestamp() *time.Time {
	if entry.Updated != nil {
		return entry.Updated
	} else if entry.Published != nil {
		return entry.Published
	} else {
		return &time.Time{}
	}
}

func (entry *Entry) CompareDate(other *Entry) int {
	return entry.Timestamp().Compare(*other.Timestamp())
}

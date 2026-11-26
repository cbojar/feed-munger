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

func (entry *Entry) CompareUpdated(other *Entry) int {
	if entry.Updated != nil && other.Updated != nil {
		return entry.Updated.Compare(*other.Updated)
	} else if entry.Updated == nil && other.Updated != nil {
		return -1
	} else if entry.Updated != nil && other.Updated == nil {
		return 1
	} else {
		return 0
	}
}

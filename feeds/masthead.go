package feeds

import (
	"time"

	c "github.com/cbojar/feed-munger/configuration"
)

type Masthead struct {
	ID           string
	Title        string
	Alternate    string
	Self         string
	Authors      []*Person
	Contributors []*Person
	Categories   []*Category
	Copyright    string
	Description  string
	Image        string
	Language     string
	Links        []*Link
	Published    *time.Time
	Updated      *time.Time
}

func MastheadFrom(options *c.Options) *Masthead {
	return &Masthead{
		ID:           options.ID(),
		Title:        options.Name(),
		Alternate:    options.Alternate(),
		Self:         options.Self(),
		Authors:      []*Person{},
		Contributors: []*Person{},
		Categories:   []*Category{},
		Links:        []*Link{},
	}
}

func (id *Masthead) FillInWith(other *Masthead) *Masthead {
	if other == nil {
		return id
	}

	return &Masthead{
		ID:           fillStringValue(id.ID, other.ID),
		Title:        fillStringValue(id.Title, other.Title),
		Alternate:    fillStringValue(id.Alternate, other.Alternate),
		Self:         fillStringValue(id.Self, other.Self),
		Authors:      other.Authors,
		Contributors: other.Contributors,
		Categories:   other.Categories,
		Copyright:    other.Copyright,
		Description:  fillStringValue(id.Description, other.Description),
		Image:        fillStringValue(id.Image, other.Image),
		Language:     fillStringValue(id.Language, other.Language),
		Links:        other.Links,
		Published:    fillTimeValue(id.Published, other.Published),
		Updated:      fillTimeValue(id.Updated, other.Updated),
	}
}

func fillStringValue(this string, other string) string {
	if this != "" {
		return this
	} else {
		return other
	}
}

func fillTimeValue(this *time.Time, other *time.Time) *time.Time {
	if this != nil {
		return this
	} else {
		return other
	}
}

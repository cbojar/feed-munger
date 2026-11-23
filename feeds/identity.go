package feeds

import (
	c "github.com/cbojar/feed-munger/configuration"
)

type Identity struct {
	id        string
	title     string
	alternate string
	self      string
}

func NewIdentity(id string, title string, alternate string, self string) *Identity {
	return &Identity{
		id:        id,
		title:     title,
		alternate: alternate,
		self:      self,
	}
}

func IdentityFrom(options *c.Options) *Identity {
	return &Identity{
		id:        options.ID(),
		title:     options.Name(),
		alternate: options.Alternate(),
		self:      options.Self(),
	}
}

func (id *Identity) ID() string {
	return id.id
}

func (id *Identity) Title() string {
	return id.title
}

func (id *Identity) Alternate() string {
	return id.alternate
}

func (id *Identity) Self() string {
	return id.self
}

func (id *Identity) FillInWith(other *Identity) *Identity {
	var idValue string
	if id.id != "" {
		idValue = id.id
	} else {
		idValue = other.ID()
	}

	var title string
	if id.title != "" {
		title = id.title
	} else {
		title = other.Title()
	}

	var alternate string
	if id.alternate != "" {
		alternate = id.alternate
	} else {
		alternate = other.Alternate()
	}

	var self string
	if id.self != "" {
		self = id.self
	} else {
		self = other.Self()
	}

	return &Identity{
		id:        idValue,
		title:     title,
		alternate: alternate,
		self:      self,
	}
}

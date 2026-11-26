package feeds

import (
	"github.com/beevik/etree"
)

type Person struct {
	Name       string
	URI        string
	Email      string
	Extensions []*etree.Element
}

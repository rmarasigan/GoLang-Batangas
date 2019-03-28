package models

import (
	"github.com/uadmin/uadmin"
)

type ImportFolder struct {
	uadmin.Model
	Path string
}

func (i ImportFolder) String() string {
	return i.Path
}

func (i *ImportFolder) Save() {
	uadmin.Save(i)
	// Some magic code
}

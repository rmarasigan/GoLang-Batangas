package models

import (
	"github.com/uadmin/uadmin"
)

type Folder struct {
	uadmin.Model
	Name     string
	Parent   *Folder
	ParentID uint
}

func (f Folder) String() string {
	return f.Name
}

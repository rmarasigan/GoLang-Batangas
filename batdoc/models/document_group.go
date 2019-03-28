package models

import (
	"github.com/uadmin/uadmin"
)

type DocumentGroup struct {
	uadmin.Model
	Group      uadmin.UserGroup
	GroupID    uint
	Document   Document
	DocumentID uint
	Read       bool
	Add        bool
	Edit       bool
	Delete     bool
}

func (DocumentGroup) HideInDashboard() bool {
	return true
}

package models

import (
	"github.com/uadmin/uadmin"
)

type DocumentUser struct {
	uadmin.Model
	User       uadmin.User
	UserID     uint
	Document   Document
	DocumentID uint
	Read       bool
	Add        bool
	Edit       bool
	Delete     bool
}

func (DocumentUser) HideInDashboard() bool {
	return true
}

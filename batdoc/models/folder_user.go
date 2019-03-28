package models

import (
	"github.com/uadmin/uadmin"
)

type FolderUser struct {
	uadmin.Model
	User     uadmin.User
	UserID   uint
	Folder   Folder
	FolderID uint
	Read     bool
	Add      bool
	Edit     bool
	Delete   bool
}

func (FolderUser) HideInDashboard() bool {
	return true
}

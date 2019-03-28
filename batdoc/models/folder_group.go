package models

import (
	"github.com/uadmin/uadmin"
)

type FolderGroup struct {
	uadmin.Model
	Group    uadmin.UserGroup
	GroupID  uint
	Folder   Folder
	FolderID uint
	Read     bool
	Add      bool
	Edit     bool
	Delete   bool
}

func (FolderGroup) HideInDashboard() bool {
	return true
}

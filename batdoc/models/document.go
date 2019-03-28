package models

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"time"
)

type Document struct {
	uadmin.Model
	Name        string
	File        string `uadmin:"file"`
	Description string `uadmin:"html"`
	RawText     string `uadmin:"list_exclude"`
	Format      Format `uadmin:"list_exclude"`
	Folder      Folder `uadmin:"filter"`
	FolderID    uint
	CreatedDate time.Time
	Channel     Channel `uadmin:"list_exclude"`
	ChannelID   uint
}

func (d *Document) Save() {
	docChange := false
	if d.ID != 0 {
		oldDoc := Document{}
		uadmin.Get(&oldDoc, "id = ?", d.ID)
		if d.File != oldDoc.File {
			docChange = true
		}
	} else {
		docChange = true
	}
	uadmin.Save(d)
	if docChange {
		ver := DocumentVersion{}
		ver.Date = time.Now()
		ver.DocumentID = d.ID
		ver.File = d.File
		ver.Number = uadmin.Count([]DocumentVersion{}, "document_id = ?", d.ID) + 1
		uadmin.Save(&ver)
	}
}

func (d Document) GetPermissions(user uadmin.User) (Read bool, Add bool, Edit bool, Delete bool) {
	if user.Admin {
		Read = true
		Add = true
		Edit = true
		Delete = true
	}
	// We will check for folder permissions first
	// Then we will check for document permissions after that
	if d.FolderID != 0 {
		folderGroup := FolderGroup{}
		uadmin.Get(&folderGroup, "group_id = ? AND folder_id = ?", user.UserGroupID, d.FolderID)
		if folderGroup.ID != 0 {
			Read = folderGroup.Read
			Add = folderGroup.Add
			Edit = folderGroup.Edit
			Delete = folderGroup.Delete
		}
		folderUser := FolderUser{}
		uadmin.Get(&folderUser, "user_id = ? AND folder_id = ?", user.ID, d.FolderID)
		if folderUser.ID != 0 {
			Read = folderUser.Read
			Add = folderUser.Add
			Edit = folderUser.Edit
			Delete = folderUser.Delete
		}
	}

	// Document Permissions
	documentGroup := DocumentGroup{}
	uadmin.Get(&documentGroup, "group_id = ? AND document_id = ?", user.UserGroupID, d.ID)
	if documentGroup.ID != 0 {
		Read = documentGroup.Read
		Add = documentGroup.Add
		Edit = documentGroup.Edit
		Delete = documentGroup.Delete
	}

	documentUser := DocumentUser{}
	uadmin.Get(&documentUser, "user_id = ? AND document_id = ?", user.ID, d.ID)
	if documentUser.ID != 0 {
		Read = documentUser.Read
		Add = documentUser.Add
		Edit = documentUser.Edit
		Delete = documentUser.Delete
	}
	return
}

func (d Document) Permissions__Form() string {
	u := uadmin.User{}
	uadmin.Get(&u, "id = ?", 2)
	r, a, e, del := d.GetPermissions(u)
	return fmt.Sprintf("Read: %v Add: %v Edit: %v Delete: %v", r, a, e, del)
}

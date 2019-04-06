package models

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"strings"
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
	CreatedBy   string
}

func (d *Document) Save() {
	docChange := false
	newDoc := false
	if d.ID != 0 {
		oldDoc := Document{}
		uadmin.Get(&oldDoc, "id = ?", d.ID)
		if d.File != oldDoc.File {
			docChange = true
		}
	} else {
		docChange = true
		newDoc = true
	}
	uadmin.Save(d)
	if docChange {
		ver := DocumentVersion{}
		ver.Date = time.Now()
		ver.DocumentID = d.ID
		ver.File = d.File
		ver.Number = uadmin.Count([]DocumentVersion{}, "document_id = ?", d.ID) + 1
		uadmin.Save(&ver)

		if newDoc {
			user := uadmin.User{}
			uadmin.Get(&user, "username = ?", d.CreatedBy)
			creator := DocumentUser{
				UserID:     user.ID,
				DocumentID: d.ID,
				Read:       true,
				Edit:       true,
				Add:        true,
				Delete:     true,
			}
			uadmin.Save(&creator)
		}
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

func (d Document) Count(a interface{}, query interface{}, args ...interface{}) int {
	Q := fmt.Sprint(query)

	if strings.Contains(Q, "user_id = ?") {
		qParts := strings.Split(Q, " AND ")
		tempArgs := []interface{}{}
		tempQuery := []string{}
		for i := range qParts {
			if qParts[i] != "user_id = ?" {
				tempArgs = append(tempArgs, args[i])
				tempQuery = append(tempQuery, qParts[i])
			}
		}
		query = strings.Join(tempQuery, " AND ")
		args = tempArgs
	}
	return uadmin.Count(a, query, args...)
}

func (d Document) AdminPage(order string, asc bool, offset int, limit int, a interface{}, query interface{}, args ...interface{}) (err error) {
	if offset < 0 {
		offset = 0
	}
	userID := uint(0)

	Q := fmt.Sprint(query)

	if strings.Contains(Q, "user_id = ?") {
		uadmin.Trail(uadmin.DEBUG, "1")
		qParts := strings.Split(Q, " AND ")
		tempArgs := []interface{}{}
		tempQuery := []string{}
		for i := range qParts {
			if qParts[i] != "user_id = ?" {
				tempArgs = append(tempArgs, args[i])
				tempQuery = append(tempQuery, qParts[i])
			} else {
				uadmin.Trail(uadmin.DEBUG, "UserID: %d", args[i])
				userID, _ = (args[i]).(uint)
			}
		}
		query = strings.Join(tempQuery, " AND ")
		args = tempArgs
	}
	if userID == 0 {
		uadmin.Trail(uadmin.DEBUG, "2")
		err = uadmin.AdminPage(order, asc, offset, limit, a, query, args...)
		return err
	}

	user := uadmin.User{}
	uadmin.Get(&user, "id = ?", userID)

	docList := []Document{}
	tempList := []Document{}
	for {
		err = uadmin.AdminPage(order, asc, offset, limit, &tempList, query, args)
		uadmin.Trail(uadmin.DEBUG, "8: offset:%d, limit:%s", offset, limit)
		if err != nil {
			uadmin.Trail(uadmin.DEBUG, "3")
			*a.(*[]Document) = docList
			return err
		}
		if len(tempList) == 0 {
			uadmin.Trail(uadmin.DEBUG, "4")

			*a.(*[]Document) = docList
			uadmin.Trail(uadmin.DEBUG, "a: %#v", a)
			return nil
		}
		for i := range tempList {
			p, _, _, _ := tempList[i].GetPermissions(user)
			if p {
				uadmin.Trail(uadmin.DEBUG, "5")
				docList = append(docList, tempList[i])
			}
			if len(docList) == limit {
				uadmin.Trail(uadmin.DEBUG, "6")
				*a.(*[]Document) = docList
				return nil
			}
		}
		offset += limit
	}
	*a.(*[]Document) = docList
	uadmin.Trail(uadmin.DEBUG, "7")
	return nil
}

func (d Document) Permissions__Form() string {
	u := uadmin.User{}
	uadmin.Get(&u, "id = ?", 2)
	r, a, e, del := d.GetPermissions(u)
	return fmt.Sprintf("Read: %v Add: %v Edit: %v Delete: %v", r, a, e, del)
}

package main

import (
	"github.com/rmarasigan/golang-batangas/batdoc/models"
	"github.com/uadmin/uadmin"
)

func main() {
	uadmin.Register(
		models.Document{},
		models.Channel{},
		models.DocumentUser{},
		models.DocumentGroup{},
		models.DocumentVersion{},
		models.Folder{},
		models.FolderUser{},
		models.FolderGroup{},
		models.ImportFolder{},
	)

	uadmin.RegisterInlines(
		models.Document{},
		map[string]string{
			"documentversion": "DocumentID",
			"documentgroup":   "DocumentID",
			"documentuser":    "DocumentID",
		},
	)

	uadmin.RegisterInlines(
		models.Folder{},
		map[string]string{
			"foldergroup": "FolderID",
			"folderuser":  "FolderID",
		},
	)

	docS := uadmin.Schema["document"]
	docS.ListModifier = DocumentListFilter
	uadmin.Schema["document"] = docS

	uadmin.SiteName = "Bat Doc"
	uadmin.StartServer()
}

func DocumentListFilter(s *uadmin.ModelSchema, u *uadmin.User) (string, []interface{}) {
	if !u.Admin {
		return "user_id = ?", []interface{}{u.ID}
	}
	return "", []interface{}{}
}

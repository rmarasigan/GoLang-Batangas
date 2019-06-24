package models

import (
	"os"
	"testing"
	"time"

	"github.com/uadmin/uadmin"
)

func TestMain(t *testing.M) {
	go func() {
		uadmin.Register(
			Document{},
			Channel{},
			DocumentUser{},
			DocumentGroup{},
			DocumentVersion{},
			Folder{},
			FolderUser{},
			FolderGroup{},
			ImportFolder{},
		)

		uadmin.RegisterInlines(
			Document{},
			map[string]string{
				"documentversion": "DocumentID",
				"documentgroup":   "DocumentID",
				"documentuser":    "DocumentID",
			},
		)

		uadmin.RegisterInlines(
			Folder{},
			map[string]string{
				"foldergroup": "FolderID",
				"folderuser":  "FolderID",
			},
		)

		//docS := uadmin.Schema["document"]
		//docS.ListModifier = DocumentListFilter
		//uadmin.Schema["document"] = docS

		uadmin.SiteName = "Bat Doc"
		uadmin.StartServer()
	}()
	time.Sleep(time.Second * 3)
	retCode := t.Run()
	os.Exit(retCode)
}

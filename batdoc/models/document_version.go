package models

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"time"
)

type DocumentVersion struct {
	uadmin.Model
	Document   Document
	DocumentID uint
	File       string `uadmin:"file"`
	Number     int    `uadmin:"help:version number"`
	Date       time.Time
	Format     Format
	RawText    string `uadmin:"list_exclude;hidden" sql:"type:text;"`
}

func (d DocumentVersion) String() string {
	return fmt.Sprint(d.Number)
}

func (DocumentVersion) HideInDashboard() bool {
	return true
}

func (d *DocumentVersion) Save() {
	newDoc := false
	uadmin.Trail(uadmin.DEBUG, "DocumentVersion.ID = %d", d.ID)
	if d.ID == 0 {
		newDoc = true
	}
	uadmin.Save(d)

	// Run Document analysis in a separate gorotine
	if newDoc {
		go func() {
			// Get Raw Text
			uadmin.Trail(uadmin.DEBUG, "Getting Raw Text")
			d.RawText = d.Format.GetRawText(d)
			uadmin.Save(d)
			uadmin.Preload(d)
			d.Document.RawText = d.RawText
			uadmin.Save(&d.Document)
		}()
	}
}

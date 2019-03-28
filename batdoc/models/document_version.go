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
}

func (d DocumentVersion) String() string {
	return fmt.Sprint(d.Number)
}

func (DocumentVersion) HideInDashboard() bool {
	return true
}

package models

import (
	"github.com/uadmin/uadmin"
)

type Channel struct {
	uadmin.Model
	Name string `uadmin:"required"`
}

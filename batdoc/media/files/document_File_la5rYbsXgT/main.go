package main

import (
	"github.com/twistedhardware/license_server/models"
	"github.com/uadmin/uadmin"
)

func main() {
	uadmin.Register(
		models.License{},
		models.Renewal{},
	)
	uadmin.Trail(uadmin.DEBUG, uadmin.GenerateBase64(10))
	uadmin.SiteName = "License Server"
	uadmin.StartServer()
}

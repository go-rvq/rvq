package main

import (
	"github.com/go-rvq/rvq/admin/example/admin"
	"github.com/go-rvq/rvq/admin/publish"
)

func main() {
	db := admin.ConnectDB()
	config := admin.NewConfig(db)
	storage := admin.PublishStorage
	publish.RunPublisher(db, storage, config.Publisher)
	select {}
}

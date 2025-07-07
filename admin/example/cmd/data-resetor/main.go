package main

import (
	"github.com/go-rvq/rvq/admin/example/admin"
)

func main() {
	db := admin.ConnectDB()
	tbs := admin.GetNonIgnoredTableNames(db)
	admin.EmptyDB(db, tbs)
	admin.InitDB(db, tbs)
	admin.ErasePublicUsersData(db)
	return
}

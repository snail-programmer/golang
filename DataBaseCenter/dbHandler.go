package DBCenter

import (
	"database/sql"

	_ "../lib/mysql"
)

var db *sql.DB

func checkDbError(err error) {
	if err != nil {
		panic(err)
	}
}
func init() {
	var err error
	db, err = sql.Open("mysql", "root:root@/schole?charset=utf8")
	checkDbError(err)
}

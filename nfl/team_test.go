package nfl

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func dbTest() {
    db, _ := sql.Open("mysql", "root:arizona9@tcp(192.168.2.101:3306)/nfl?autocommit=true")
    defer db.Close()
    result, _ := db.Query(db.Prepare("select 1"))
    result.
}

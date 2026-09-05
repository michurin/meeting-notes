package common

import (
	"database/sql"
	"fmt"
)

func Dump(db *sql.DB) {
	rows := Must2(db.Query("select name from city"))
	defer rows.Close()
	fmt.Println("--- dump ---")
	for rows.Next() {
		var name string
		Must(rows.Scan(&name))
		fmt.Println(name)
	}
	Must(rows.Err())
}

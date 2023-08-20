package db

import (
	"elichika/serverdb"

	"fmt"
)

func Init(args []string ) {
	// wipe the database if necessary, then insert the relevant table
	if len(args) == 0 {
		serverdb.InitTables(false)
	} else {
		if args[0] == "overwrite" {
			serverdb.InitTables(true)
		} else {
			fmt.Println("Invalid params:", args)
			fmt.Println("Use [overwrite] to overwrite existing data, or nothing to just init new table")
		}
	}
}

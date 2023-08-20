package main

import (
	"elichika/router"
	
	"os"
	"fmt"
	"elichika/cli/db"

	"github.com/gin-gonic/gin"
)

func make(args []string) {
	server := args[0]
	if (server != "jp") && (server != "gl") {
		fmt.Println("Server must be \"jp\" or \"gl\", not", server)
	}
	accountType := args[1]
	if (accountType != "new") && (accountType != "json") {
		fmt.Println("Account type must be \"new\" or \"gl\", json", accountType)
	}

	db.Init([]string{"overwrite"})
	db.Account(args)
	db.Gacha([]string{"init"})
	db.Gacha([]string{"insert", "gacha_insert.json"})
}

func cli() {
	switch os.Args[1] {
	case "init":
		db.Init(os.Args[2:])
	case "account":
		db.Account(os.Args[2:])
	case "gacha":
		db.Gacha(os.Args[2:])
	case "make": // easy import
		make(os.Args[2:])
	default:
		fmt.Println("Invalid params:", os.Args)
		return
	}
	return
}


func main() {

	if len(os.Args) > 1 {
		cli()
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	router.Router(r)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

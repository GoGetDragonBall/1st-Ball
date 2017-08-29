package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var configDir = "conf/"

type DatabaseConf struct {
	DriverName, DataSourceName string
}

func main() {
	dbConfigFile, err := os.Open(configDir + "database.json")
	if err != nil {
		panic(err)
	}

	var dbConf DatabaseConf
	dbConfigDecoder := json.NewDecoder(dbConfigFile)
	err = dbConfigDecoder.Decode(&dbConf)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(dbConf.DriverName, dbConf.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var name string
	var nickname string
	var email string
	row := db.QueryRow("SELECT name, nickname, email FROM users")
	if err = row.Scan(&name, &nickname, &email); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := []byte("name: " + name + " nickname: " + nickname + " email: " + email)
		w.Write(body)
	})

	http.ListenAndServe(":9000", nil)
}

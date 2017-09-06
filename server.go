package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		log.Println("in")
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users/signup", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":

			//password := r.FormValue("password")
			name := r.FormValue("name")
			nickname := r.FormValue("nickname")
			email := r.FormValue("email")

			result, err := db.Exec("INSERT INTO users (name,nickname,email) values(?,?,?)", name, nickname, email)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)

		case "GET":
			content, err := ioutil.ReadFile("html/user_insert.html")
			if err != nil {
				log.Fatal(err)
			}
			w.Write(content)
		}
	})

	http.ListenAndServe(":9000", nil)
}

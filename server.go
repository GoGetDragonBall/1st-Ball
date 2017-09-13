package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"golang.org/x/crypto/bcrypt"
)

// FIXME: 한글 안 들어감. ㅅㅂㅅㅂㅅㅂ
type User struct {
	ID uint
	Name string
	Nickname string
	Email string `gorm:"type:varchar(100)"`
	UserAuthorizationKeys []UserAuthorizationKey
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserAuthorizationKey struct {
	ID uint
	AuthorizationKey string `gorm:"type:varchar(300)"`
	TypeId uint
	UserId uint
	UserAuthorizationKeyTypes []UserAuthorizationKeyType `gorm:"ForeignKey:ID;AssociationForeignKey:TypeId"`
	Users []User `gorm:"ForeignKey:ID;AssociationForeignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserAuthorizationKeyType struct {
	ID uint
	Name string `gorm:"type:varchar(100)"`
	UserAuthorizationKeys []UserAuthorizationKey
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

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

	db, err := gorm.Open(dbConf.DriverName, dbConf.DataSourceName)
	if err != nil {
		log.Println("in")
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users/signup", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":
			name := r.FormValue("name")
			password := r.FormValue("password")
			nickname := r.FormValue("nickname")
			email := r.FormValue("email")

			bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}

			var authorizationKeyBcryptType UserAuthorizationKeyType
			db.Where("name = ?", "bcrypt").Find(&authorizationKeyBcryptType)

			user := User{
				Name: name,
				Nickname: nickname,
				Email: email,
				UserAuthorizationKeys: []UserAuthorizationKey{
					{AuthorizationKey: string(bcryptedPassword), TypeId: authorizationKeyBcryptType.ID},
				},
			}

			userCreateErrs := db.Create(&user).GetErrors()
			if len(userCreateErrs) != 0 {
				for err := range userCreateErrs {
					log.Println(err)
				}

				w.WriteHeader(500)
				w.Write([]byte("{\"message\":\"WTF\"}"))
			} else {
				w.WriteHeader(200)
				w.Write([]byte("{\"message\":\"Good\"}"))
			}

			return
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

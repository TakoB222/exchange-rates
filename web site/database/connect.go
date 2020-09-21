package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var Db *sql.DB

func init(){
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	dbb, err := sql.Open("postgres", dbinfo)
	if err != nil{
		fmt.Println("db open error: ", err)
	}
	Db = dbb

}

func GetDB()*sql.DB{
	return Db
}

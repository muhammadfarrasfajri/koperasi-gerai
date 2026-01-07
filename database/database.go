package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectMySQL() {
	username := os.Getenv("USERNAME_DATABASE")
	database := os.Getenv("DATABASE")
	database_name := os.Getenv("DATABASE_NAME")
	
	var err error
	dsn := fmt.Sprintf("%s:@tcp(127.0.0.1:3306)/%s?parseTime=true&loc=Local", username, database_name)
	DB, err = sql.Open(database,dsn)
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable: ", err)
	}

	log.Println("MySQL connected.")
}

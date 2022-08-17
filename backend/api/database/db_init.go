package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

// variable that points to the postgres connection, allows web server to interact with database
var Db *sql.DB

func init() {
	//get db credentials from env. file
	err := godotenv.Load("postgres.env")
	if err != nil {
		panic(err)
	}

	var (
		host     = os.Getenv("host")
		port, _  = strconv.Atoi(os.Getenv("port"))
		user     = os.Getenv("user")
		password = os.Getenv("password")
		dbname   = os.Getenv("dbname")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//connect to db server
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

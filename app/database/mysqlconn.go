package db

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	UserName string
	Password string
	Host     string
	Port     int
	DbName   string
}

var Db *sql.DB

func InitDB(C DbConfig) *sql.DB {
	// Initialize the database connection
	var err error
	dataSourceName := C.UserName + ":" + C.Password + "@tcp(" + C.Host + ":" + strconv.Itoa(C.Port) + ")/" + C.DbName
	Db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err)
	}

	// Test the database connection
	if err := Db.Ping(); err != nil {
		panic(err)
	}

	createTables()
	insertSamplePage()
	return Db
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS pages (
			id smallint(8) unsigned NOT NULL auto_increment,
			link varchar(2000) NULL,
			content varchar(100000) NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB;`,
	}
	for _, q := range queries {
		_, err := Db.Exec(q)
		if err != nil {
			panic(err.Error())
		}
	}
}

func insertSamplePage() {
	result, err := Db.Exec(
		"REPLACE INTO pages (id,link,content) VALUES(?,?,?)",
		1,
		"cms-page",
		"Hello World! <div><img src='/assets/media/images/pic.png'></div><a href='/'>home</a>",
	)
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	if err != nil || lastid < 1 {
		panic(err)
	}
}

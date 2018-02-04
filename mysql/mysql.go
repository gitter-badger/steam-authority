package mysql

import (
	"os"

	"github.com/Jleagle/go-helpers/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mysqlConnection *sqlx.DB

func getDB() (conn *sqlx.DB, err error) {

	if mysqlConnection == nil {

		db, err := sqlx.Connect("mysql", os.Getenv("STEAM_SQL_DSN")+"?parseTime=true")
		if err != nil {
			logger.Error(err)
			return db, err
		}

		mysqlConnection = db
	}

	return mysqlConnection, nil
}

func CountTable(table string) (count uint, err error) {

	db, err := getDB()
	if err != nil {
		return count, err
	}

	err = db.Get(&count, "SELECT count(*) as id FROM "+table)
	if err != nil {
		return count, err
	}

	return count, nil
}

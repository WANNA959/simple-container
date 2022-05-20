package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriverName = "sqlite3"
	dbPath       = "/tmp/sc.db"
)

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

func InitSqlite() (err error) {
	db, err = sql.Open(dbDriverName, dbPath)
	if err != nil {
		return err
	}
	err = createTable()
	if err != nil {
		return err
	}
	return nil
}

func createTable() error {
	// create table network_mgr
	sql := `create table if not exists "container_mgr" (
		"id" integer primary key autoincrement,
		"pid" text not null,
		"veth" integer not null,
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    	"update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
	)`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	// trigger for update_time
	sql = `
	CREATE TRIGGER if not exists update_time_trigger UPDATE OF id,pid,veth,create_time ON container_mgr
	BEGIN
	  UPDATE network_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
	END
	`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// create table network_mgr
	sql = `create table if not exists "network_mgr" (
		"id" integer primary key autoincrement,
		"subnet" text not null,
		"bind_ip" text not null default "",
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    	"update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
	)`

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	// trigger for update_time
	sql = `
	CREATE TRIGGER if not exists update_time_trigger2 UPDATE OF id,subnet,bind_ip,create_time ON network_mgr
	BEGIN
	  UPDATE network_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
	END
	`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	return err
}

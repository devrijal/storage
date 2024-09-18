package storage

import (
	"log"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestMySQL(t *testing.T) {

	options := make(map[string]string)
	options["parseTime"] = "True"

	dsn := &RelationalDatabase{
		Type:     MySQL,
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "123",
		Database: "tests",
		Option:   options,
		Config:   gorm.Config{},
	}

	db, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}

	type Version string
	var version Version

	q := db.Debug().Raw("select @@version as version").Scan(&version)

	if q.Error != nil {
		t.Fatal(q.Error)
	}

	log.Printf("Version: %v", version)
}

func TestMariaDB(t *testing.T) {

	options := make(map[string]string)
	options["parseTime"] = "True"

	dsn := &RelationalDatabase{
		Type:     MariaDB,
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "123",
		Database: "tests",
		Option:   options,
		Config:   gorm.Config{},
	}

	db, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}

	type Version string
	var version Version

	q := db.Debug().Raw("select @@version as version").Scan(&version)

	if q.Error != nil {
		t.Fatal(q.Error)
	}

	log.Printf("Version: %v", version)
}

func TestPostgreSQL(t *testing.T) {

	options := make(map[string]string)
	options["sslmode"] = "disable"
	options["TimeZone"] = "Asia/Jakarta"

	dsn := &RelationalDatabase{
		Type:     PostgreSQL,
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "123",
		Database: "tests",
		Option:   options,
		Config:   gorm.Config{},
	}

	db, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}

	type Version string
	var version Version

	q := db.Debug().Raw("select version()").Scan(&version)

	if q.Error != nil {
		t.Fatal(q.Error)
	}

	log.Printf("Version: %v", version)
}

type SalesOrderHeader struct {
	SONo string
}

func TestSQLServer(t *testing.T) {
	dsn := &RelationalDatabase{
		Type:     SQLServer,
		Host:     "127.0.0.1",
		Port:     "1433",
		User:     "root",
		Password: "123",
		Database: "tests",
		Config: gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
				NoLowerCase:   true, // skip the snake_casing of names
			},
			DisableForeignKeyConstraintWhenMigrating: true},
	}

	db, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}

	type Version string
	var version Version

	q := db.Debug().Raw("select @@version").Scan(&version)

	if q.Error != nil {
		t.Fatal(q.Error)
	}

	log.Printf("Version: %v", version)
}

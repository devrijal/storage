package storage

import (
	"testing"

	"gorm.io/gorm"
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

	db.Debug().Raw("Select @@version")
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

	db.Debug().Raw("Select @@version")
}

func TestPostgreSQL(t *testing.T) {

	options := make(map[string]string)
	options["sslmode"] = "disable"
	options["TimeZone"] = "Asia/Jakarta"

	dsn := &RelationalDatabase{
		Type:     PostgreSQL,
		Host:     "127.0.0.1",
		Port:     "5433",
		User:     "postgres",
		Password: "123",
		Database: "postgres",
		Option:   options,
		Config:   gorm.Config{},
	}

	_, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSQLServer(t *testing.T) {
	dsn := &RelationalDatabase{
		Type:     SQLServer,
		Host:     "127.0.0.1",
		Port:     "1433",
		User:     "test",
		Password: "test",
		Database: "tests",
		Config:   gorm.Config{},
	}

	_, err := Relational(dsn)

	if err != nil {
		t.Fatalf(err.Error())
	}
}

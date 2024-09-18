package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	connections map[string]*gorm.DB
)

type RDBMS int

const (
	MySQL RDBMS = iota
	MariaDB
	PostgreSQL
	SQLServer
)

func (db RDBMS) String() string {
	return [...]string{"MySQL", "MariaDB", "PostgreSQL", "SQLServer"}[db]
}

func (db RDBMS) Index() int {
	return int(db)
}

type RelationalDatabase struct {
	Type     RDBMS
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Option   map[string]string
	Config   gorm.Config
}

func (db *RelationalDatabase) GetType() RDBMS {
	return db.Type
}
func (db *RelationalDatabase) GetConnectionString() string {
	var dsn string

	switch db.Type {
	case MySQL:
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database,
		)

		for key, val := range db.Option {

			dsn = fmt.Sprintf("%s&%s=%s", dsn, string(key), val)
		}
	case MariaDB:
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database,
		)

		for key, val := range db.Option {

			dsn = fmt.Sprintf("%s&%s=%s", dsn, string(key), val)
		}

	case PostgreSQL:
		str := "host=%s user=%s password=%s dbname=%s port=%s"

		dsn = fmt.Sprintf(str, db.Host, db.User, db.Password, db.Database, db.Port)

		for key, val := range db.Option {

			dsn = fmt.Sprintf("%s %s=%s", dsn, string(key), val)
		}
	case SQLServer:
		dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database,
		)
	}

	return dsn
}

func (db *RelationalDatabase) GetConfig() *gorm.Config {
	return &db.Config
}

func (r *RelationalDatabase) GetRetryInterval() time.Duration {
	return time.Duration(3)
}

func (db *RelationalDatabase) GetDialector() gorm.Dialector {
	var dialector gorm.Dialector

	switch db.Type {
	case MySQL:
		dialector = mysql.Open(db.GetConnectionString())
	case MariaDB:
		dialector = mysql.Open(db.GetConnectionString())
	case PostgreSQL:
		dialector = postgres.Open(db.GetConnectionString())
	case SQLServer:
		dialector = sqlserver.Open(db.GetConnectionString())
	}

	return dialector
}

type Manager interface {
	GetConnectionString() string
}

type RelationalDatabaseManager interface {
	Manager
	GetType() RDBMS // Get storage type
	GetRetryInterval() time.Duration
	GetConfig() *gorm.Config
	GetDialector() gorm.Dialector
}

func connectRdbms(s RelationalDatabaseManager) (*gorm.DB, error) {

	db, err := gorm.Open(s.GetDialector(), s.GetConfig())

	if err != nil {
		return nil, err
	}

	return db, nil
}

func reconnectRdbms(s RelationalDatabaseManager) (*gorm.DB, error) {
	interval := s.GetRetryInterval()

	for {
		db, err := connectRdbms(s)

		if err == nil {
			return db, nil
		}

		time.Sleep(interval)
	}
}

func Relational(s RelationalDatabaseManager) (db *gorm.DB, err error) {

	dsn := s.GetConnectionString()

	if len(connections) == 0 {
		connections = make(map[string]*gorm.DB)
	}

	if connections[dsn] == nil {

		db, err := reconnectRdbms(s)

		if err != nil {
			return db, err
		}

		connections[dsn] = db
	}

	return connections[dsn], err
}

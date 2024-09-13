package storage

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	connections map[string]*gorm.DB
	once        sync.Once
)

const (
	MySQL = iota
	MariaDB
	PostgreSQL
	SQLServer
)

type RelationalDatabase struct {
	Type     int
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Option   map[string]string
	Config   gorm.Config
}

func (db *RelationalDatabase) GetType() int {
	return db.Type
}
func (db *RelationalDatabase) GetConnectionString() string {
	var dsn string

	switch db.Type {
	case 0:
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
	case 1:
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
	case 2:
		str := "host=%s user=%s password=%s dbname=%s port=%s"

		dsn = fmt.Sprintf(str, db.Host, db.User, db.Password, db.Database, db.Port)

		for key, val := range db.Option {

			dsn = fmt.Sprintf("%s %s=%s", dsn, string(key), val)
		}
	case 3:
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
	case 0:
		dialector = mysql.Open(db.GetConnectionString())
	case 1:
		dialector = mysql.Open(db.GetConnectionString())
	case 2:
		dialector = postgres.Open(db.GetConnectionString())
	case 3:
		sqlserver.Open(db.GetConnectionString())
	}

	return dialector
}

type Manager interface {
	GetConnectionString() string
}

type RelationalDatabaseManager interface {
	Manager
	GetType() int // Get storage type
	GetRetryInterval() time.Duration
	GetConfig() *gorm.Config
	GetDialector() gorm.Dialector
}

func connectRdbms(s RelationalDatabaseManager) (*gorm.DB, error) {
	return gorm.Open(s.GetDialector(), s.GetConfig())
}

func reconnectRdbms(s RelationalDatabaseManager) (*gorm.DB, error) {
	interval := s.GetRetryInterval()

	for {
		db, err := connectRdbms(s)

		if err == nil {
			return db, nil
		}

		log.Printf("Cannot connect to %v: %v", s.GetType(), err)
		time.Sleep(interval)
	}
}

func Relational(s RelationalDatabaseManager) (db *gorm.DB, err error) {

	dsn := s.GetConnectionString()

	if len(connections) == 0 {
		connections = make(map[string]*gorm.DB)
	}

	if connections[dsn] == nil {
		once.Do(func() {
			db, err := reconnectRdbms(s)

			if err != nil {
				log.Fatal(err)
			}

			connections[dsn] = db
		})
	}

	return connections[dsn], err
}

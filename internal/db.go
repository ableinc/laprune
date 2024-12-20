package internal

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postgresDBInstance *sql.DB
var mongoDBInstance *mongo.Client
var sqliteDBInstance *sql.DB
var mssqlDBInstance *sql.DB
var mysqlDBInstance *sql.DB

var MongodbCtx context.Context
var lock *sync.Mutex = &sync.Mutex{}

type Vendor string

const (
	MONGODB  Vendor = "mongodb"
	POSTGRES Vendor = "postgres"
	SQLITE   Vendor = "sqlite"
	MSSQL    Vendor = "mssql"
	MYSQL Vendor = "mysql"
)

func Connection(vendor Vendor, uri string) (*sql.DB, *mongo.Client) {
	lock.Lock()
	defer lock.Unlock()
	switch vendor {
	case MONGODB:
		if mongoDBInstance != nil {
			return nil, mongoDBInstance
		}
		mongoDBInstance = startMongodbConnection(uri)
		return nil, mongoDBInstance
	case POSTGRES:
		if postgresDBInstance != nil {
			return postgresDBInstance, nil
		}
		postgresDBInstance = startPostgresConnection(uri)
		return postgresDBInstance, nil
	case SQLITE:
		if sqliteDBInstance != nil {
			return sqliteDBInstance, nil
		}
		sqliteDBInstance = startSQLiteConnection(uri)
		return sqliteDBInstance, nil
	case MSSQL:
		if mssqlDBInstance != nil {
			return mssqlDBInstance, nil
		}
		mssqlDBInstance = startMSSQLConnection(uri)
		return mssqlDBInstance, nil
	case MYSQL:
		if mysqlDBInstance != nil {
			return mysqlDBInstance, nil
		}
		mysqlDBInstance = startMysqlConnection(uri)
		return mysqlDBInstance, nil
	default:
		panic(errors.New("Invalid vendor provided"))
	}
}

func startMongodbConnection(uri string) *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	MongodbCtx = ctx
	return client
}

func startPostgresConnection(uri string) *sql.DB {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	// Return pointer to db connection
	return db
}

func startSQLiteConnection(uri string) *sql.DB {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		panic(err)
	}
	return db
}

func startMSSQLConnection(uri string) *sql.DB {
	// Create connection pool
	db, err := sql.Open("sqlserver", uri)
	if err != nil {
		panic(err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func startMysqlConnection(uri string) *sql.DB {
	// Create connection pool
	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

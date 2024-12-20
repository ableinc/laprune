package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	internal "laprune/internal"

	"go.mongodb.org/mongo-driver/bson"
)

type Environment string

const (
	TEST        Environment = "test"
	DEVELOPMENT Environment = "development"
	STAGING     Environment = "staging"
	PRODUCTION  Environment = "production"
)

type Database struct {
	Name       string          `json:"name"`
	Vendor     internal.Vendor `json:"vendor"`
	Uri        string          `json:"uri"`
	Database   string          `json:"database"`
	Collection string          `json:"collection"`
}

type DatabaseConfig struct {
	Environment Environment `json:"environment"`
	Databases   []Database  `json:"databases"`
}

type SqlQueries struct {
	Name       string
	Statements []string
}

func main() {
	// Open json database file
	file, err := os.Open("./db.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Decode the JSON data into a struct
	var databases DatabaseConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&databases)
	if err != nil {
		panic(err)
	}
	// Start process
	sqlQueries := readSqlSchema()
	executeSqlSchema(databases, sqlQueries)
}

func readSqlSchema() map[string][]string {
	file, err := os.Open("./config/queries.sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var result map[string][]string = make(map[string][]string)
	// Create a new scanner
	scanner := bufio.NewScanner(file)
	// Read file line by line
	for scanner.Scan() {
		line := scanner.Text()
		dbNameAndSqlStatement := strings.Split(line, "=")
		dbName, sqlStatement := dbNameAndSqlStatement[0], dbNameAndSqlStatement[1]
		_, ok := result[dbName]
		if !ok {
			result[dbName] = []string{strings.TrimSpace(sqlStatement)}
		} else {
			result[dbName] = append(result[dbName], strings.TrimSpace(sqlStatement))
		}
	}
	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func executeSqlSchema(databaseConfig DatabaseConfig, sqlQueries map[string][]string) {
	fmt.Fprintf(os.Stdout, "Environment: %s\n", databaseConfig.Environment)
	for _, database := range databaseConfig.Databases {
		fmt.Fprintf(os.Stdout, "Executing statement on %s:\n", database.Name)
		sqlStatements, ok := sqlQueries[database.Name]
		if !ok {
			fmt.Fprintf(os.Stderr, "No statements found for %s", database.Name)
			continue
		}
		for _, sqlStatement := range sqlStatements {
			fmt.Fprintf(os.Stdout, "%s\n", sqlStatement)
			if sqlStatement == "" {
				continue
			}
			// Execute statements
			genericSqlClient, mongoClient := internal.Connection(database.Vendor, database.Uri)
			if mongoClient != nil {
				collection := mongoClient.Database(database.Database).Collection(database.Collection)
				// Define a filter for the delete operation
				filter := bson.M{"email": bson.M{"$regex": sqlStatement}}
				// Execute the delete operation
				_, err := collection.DeleteMany(internal.MongodbCtx, filter)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
			} else {
				_, err := genericSqlClient.Exec(sqlStatement)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"cwgo_test/biz/router"  // Assume this is the generated router package
	"cwgo_test/biz/service" // Assuming this is where your services are defined
	"github.com/cloudwego/hertz/pkg/app/server"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func main() {
	// Database connection string
	dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local"

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}

	// Create an instance of Hertz server
	h := server.Default()

	// Assuming you have a service layer setup
	svc := service.NewService(db) // You need to implement this function

	// Register routes using the generated router
	router.RegisterRoutes(h, svc)

	// Start the server
	if err := h.Spin(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
		os.Exit(1)
	}

}

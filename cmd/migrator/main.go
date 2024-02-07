package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	var migrationsPath string
	var up, down bool

	userName := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbName := os.Getenv("DB_NAME")

	storagePath := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", userName, password, host, port, dbName)

	flag.StringVar(&storagePath, "storage-path", "", "database connection URL")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.BoolVar(&up, "up", false, "apply all pending migrations")
	flag.BoolVar(&down, "down", false, "revert the last applied migration")
	flag.Parse()

	if storagePath == "postgres://:@:/" || migrationsPath == "" {
		fmt.Println("storage-path and migrations-path are required")
		os.Exit(1)
	}

	// Set environment variables for dbmate
	os.Setenv("DATABASE_URL", storagePath)
	os.Setenv("DBMATE_DIR", migrationsPath)

	// Construct the dbmate command based on the flags
	var dbmateCommand *exec.Cmd
	if up {
		dbmateCommand = exec.Command("dbmate", "up")
	} else if down {
		dbmateCommand = exec.Command("dbmate", "down")
	} else {
		fmt.Println("Please specify either -up or -down to apply or revert migrations.")
		os.Exit(1)
	}

	// Run the dbmate command
	m, err := dbmateCommand.CombinedOutput()
	if err != nil {
		fmt.Printf("dbmate failed: %s\n", err)
		fmt.Printf("Output: %s\n", strings.TrimSpace(string(m)))
		os.Exit(1)
	}

	fmt.Println("dbmate operation successful")
	fmt.Println(strings.TrimSpace(string(m)))
}

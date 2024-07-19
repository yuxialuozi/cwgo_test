package codegen

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func GenerateMain() {
	// Define the data for the template
	data := struct {
		DSN string
	}{
		DSN: "username:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local",
	}

	// Parse the template
	tmpl, err := template.New("main").Parse(mainContent)
	if err != nil {
		log.Fatalf("Error parsing main template: %v", err)
	}

	dir := "test"
	fileName := "main.go"
	path := fmt.Sprintf("%s/%s", dir, fileName)

	// Create the directory if it doesn't exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Create or truncate the file
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Execute the template with the data
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("Error executing main template: %v", err)
	}

	log.Println("main.go generated successfully")
}

package codegen

import (
	"fmt"
	"os"
	"text/template"
)

func GenerateCode(queries []SQLQueryInfo, table TableInfo) {

	generateDb()

	generateDBQuery(queries, table)

}

func generateDBQuery(queries []SQLQueryInfo, table TableInfo) {

	// Create the template
	tmpl := template.Must(template.New("queryTemplate").Parse(templateContent))

	// Create the output file
	file, err := os.Create("biz/db/query/query_gen.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Execute the template with the queries slice and table info
	data := query{
		Queries: queries,
		Table:   table,
	}

	fmt.Println(data)

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Template successfully generated to biz/db/query/query_gen.go")
}

func generateDb() {
	dir := "biz/db/query"
	fileName := "db.go"
	path := fmt.Sprintf("%s/%s", dir, fileName)

	// Create the directory if it doesn't exist
	err := os.MkdirAll(dir, os.ModePerm)
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

	// Write the template content to the file
	tmpl := template.Must(template.New("queryTemplate").Parse(queryDbTemplate))
	err = tmpl.Execute(file, nil)
	if err != nil {
		fmt.Printf("Error writing template to file: %v\n", err)
	}

}

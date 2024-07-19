package main

import (
	"cwgo_test/anylize"
	"cwgo_test/codegen"
	"io/ioutil"
	"log"
)

// readFile 读取 SQL 文件内容
func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	return string(data)
}

func main() {
	filename := "sql/generated_file.sql" // SQL 文件的路径
	sql := readFile(filename)
	tableInfo := codegen.ParseSQLToStruct(sql)

	filename = "sql/user_queris.sql"
	sql = readFile(filename)
	queryInfo := codegen.ParseSQLQueryToStruct(sql)

	// Database connection string
	dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	queryName := "test"
	anylize.SQLExecutionShow(queryInfo, dsn, queryName)

	codegen.GenerateModel(tableInfo)

	codegen.GenerateCode(queryInfo, tableInfo)

	codegen.GenerateServer(queryInfo, tableInfo)

	codegen.GenerateRouter(queryInfo, tableInfo)

	codegen.GenerateMain()

}

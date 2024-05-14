package main

import (
	"fmt"
	"github.com/marianogappa/sqlparser"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Struct struct {
	Name   string
	Fields []Field
}

func main() {
	thriftIDL := `
struct User {
		1: i64 id (json:"id", db:"id");
		2: string name (json:"name", db:"name");
		3: string email (json:"email", db:"email");
		4: i64 created_at (json:"created_at", db:"created_at");
		5: i64 updated_at (json:"updated_at", db:"updated_at");
	}

struct student {
		1: i64 id (json:"id", db:"id");
		2: string name (json:"name", db:"name");
		3: string email (json:"email", db:"email");
		4: i64 created_at (json:"created_at", db:"created_at");
		5: i64 updated_at (json:"updated_at", db:"updated_at");
	}
`

	structs := parseThriftIDL(thriftIDL)

	//// 创建一个静态解析出来的结构体列表
	//structs := []Struct{
	//	{
	//		Name: "User",
	//		Fields: []Field{
	//			{Name: "Id", Type: "int64", Tag: "`json:\"id\" db:\"id\"`"},
	//			{Name: "Name", Type: "string", Tag: "`json:\"name\" db:\"name\"`"},
	//			{Name: "Email", Type: "string", Tag: "`json:\"email\" db:\"email\"`"},
	//			{Name: "CreatedAt", Type: "int64", Tag: "`json:\"created_at\" db:\"created_at\"`"},
	//			{Name: "UpdatedAt", Type: "int64", Tag: "`json:\"updated_at\" db:\"updated_at\"`"},
	//		},
	//	},
	//}

	// 创建模型并保存到文件
	for _, s := range structs {
		createModelFile(s)
		sqlFileContent := generateCRUDSQL(s)

		// 将 SQL 文件内容写入文件
		filePath := "model/test.sql"
		err := ioutil.WriteFile(filePath, []byte(sqlFileContent), 0644)
		if err != nil {
			fmt.Printf("Failed to write SQL file: %v\n", err)
			return
		}

		fmt.Printf("SQL file saved successfully: %s\n", filePath)
	}

}

// 解析idl文件
func parseThriftIDL(thriftIDL string) []Struct {
	var structs []Struct
	var currentStruct Struct

	lines := strings.Split(thriftIDL, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "struct ") {
			if currentStruct.Name != "" {
				structs = append(structs, currentStruct)
			}
			currentStruct = Struct{Name: strings.TrimSpace(line[len("struct "):strings.Index(line, " {")])}
		} else if strings.Contains(line, ":") {
			parts := strings.Fields(line)

			fieldType := parseType(parts[1])
			fieldName := toCamelCase(parts[2])
			str := parseTag(parts[3:]...)

			fieldTag := strings.ReplaceAll(str, "(", "")
			fieldTag = strings.ReplaceAll(fieldTag, ")", "")
			fieldTag = strings.ReplaceAll(fieldTag, ";", "")
			fieldTag = strings.ReplaceAll(fieldTag, ",", "")

			currentStruct.Fields = append(currentStruct.Fields, Field{Name: fieldName, Type: fieldType, Tag: fieldTag})
		}

	}

	if currentStruct.Name != "" {
		structs = append(structs, currentStruct)
	}

	return structs
}

// 解析type
func parseType(thriftType string) string {
	switch thriftType {
	case "i64":
		return "int64"
	case "string":
		return "string"
	// 其他类型处理
	default:
		return "string"
	}
}

// 解析tag
func parseTag(tagParts ...string) string {
	// 构建标签字符串
	var builder strings.Builder
	builder.WriteRune('`')
	for _, part := range tagParts {
		builder.WriteString(part) // 添加键值对
		builder.WriteRune(' ')    // 分隔符
	}
	builder.WriteString("`")

	return builder.String()
}

// 创建模型文件
func createModelFile(s Struct) {
	fileName := fmt.Sprintf("model/%s.go", s.Name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file %s: %s\n", fileName, err)
		return
	}
	defer file.Close()

	// 写入包名和导入语句
	file.WriteString("package model\n\n")

	// 写入结构体定义
	file.WriteString(fmt.Sprintf("type %s struct {\n", s.Name))
	for _, field := range s.Fields {
		file.WriteString(fmt.Sprintf("\t%s %s %s\n", field.Name, field.Type, field.Tag))
	}
	file.WriteString("}\n")

	fmt.Printf("Model created successfully: %s\n", fileName)
}

// 大驼峰命名
func toCamelCase(s string) string {
	var result strings.Builder
	words := strings.Fields(s)

	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

// 生成 CRUD SQL 文件内容
func generateCRUDSQL(s Struct) string {
	crudCode := ""

	for _, field := range s.Fields {
		// 生成 CREATE 语句
		createSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s %s);\n", s.Name, field.Name, field.Type)

		// 生成 READ 语句
		readSQL := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?;\n", s.Name, field.Name)

		// 生成 UPDATE 语句
		updateSQL := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?;\n", s.Name, field.Name, field.Name)

		// 生成 DELETE 语句 //这里假设错了
		deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;\n", s.Name, field.Name)

		// 拼接成完整的 CRUD 代码
		crudCode += fmt.Sprintf("%s\n%s\n%s\n%s\n\n", createSQL, readSQL, updateSQL, deleteSQL)
	}

	return crudCode
}

func sqlparse(sql string) {
	query, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+#v", query)
}

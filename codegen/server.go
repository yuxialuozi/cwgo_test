package codegen

import (
	"fmt"
	"os"
	"text/template"
)

func GenerateServer(queries []SQLQueryInfo, table TableInfo) {
	// Create the template
	tmpl, err := template.New("serverTemplate").Parse(serverContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// 定义目录和文件路径
	dir := "biz/service"
	fileName := fmt.Sprintf("%s_service.go", table.TableName)
	filePath := fmt.Sprintf("%s/%s", dir, fileName)

	// 创建目录（如果它不存在）
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// 创建或覆盖文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a data structure to pass to the template
	data := struct {
		StructName string
		Queries    []SQLQueryInfo
	}{
		StructName: table.TableName + "Service",
		Queries:    queries,
	}

	// Execute the template
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Service file generated successfully:", filePath)
}

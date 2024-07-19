package codegen

import (
	"fmt"
	"os"
	"text/template"
)

func GenerateRouter(queries []SQLQueryInfo, table TableInfo) {
	// 创建模板
	tmpl, err := template.New("routerContent").Parse(routerContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// 定义文件路径
	dir := fmt.Sprintf("biz/router")
	fileName := fmt.Sprintf("%s_router.go", table.TableName)
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

	// 准备模板数据
	data := struct {
		ServiceImportPath string
		ServiceStructName string
		Queries           []SQLQueryInfo
	}{
		ServiceImportPath: "cwgo_test/biz/service",
		ServiceStructName: table.TableName + "Service",
		Queries:           queries,
	}

	// 执行模板，并写入文件
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Router file generated successfully:", filePath)
}

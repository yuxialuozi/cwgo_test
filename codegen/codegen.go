package codegen

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TableInfo 包含表结构信息和AST
type TableInfo struct {
	TableName   string
	PackageName string
	Columns     []ColumnInfo
	AST         *sqlparser.DDL
}

// ColumnInfo 描述表中的一个字段
type ColumnInfo struct {
	FieldName string
	GoType    string
	JsonTag   string
	QueryTag  string
	BodyTag   string
}

// parseSQLToStruct 解析 SQL CREATE TABLE 语句并生成 TableInfo 结构体
func ParseSQLToStruct(sql string) TableInfo {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatalf("failed to parse SQL: %v", err)
	}

	createStmt, ok := stmt.(*sqlparser.DDL)
	if !ok || createStmt.Action != sqlparser.CreateStr {
		log.Fatalf("not a CREATE TABLE statement")
	}

	// Convert table name to lowercase for package name
	packageName := strings.ToLower(createStmt.NewName.Name.String())

	tableInfo := TableInfo{
		TableName:   strings.Title(createStmt.NewName.Name.String()),
		PackageName: packageName,
		AST:         createStmt,
	}

	for _, col := range createStmt.TableSpec.Columns {
		goType := sqlTypeToGoType(col.Type.Type)
		fieldName := strings.Title(col.Name.String()) // Capitalize the field name for Go struct
		columnInfo := ColumnInfo{
			FieldName: fieldName,
			GoType:    goType,
			JsonTag:   strings.ToLower(col.Name.String()),
			QueryTag:  strings.ToLower(col.Name.String()),
			BodyTag:   strings.ToLower(col.Name.String()),
		}
		tableInfo.Columns = append(tableInfo.Columns, columnInfo)
	}
	return tableInfo
}

// sqlTypeToGoType 将 SQL 类型转换为 Go 类型
func sqlTypeToGoType(sqlType string) string {
	switch strings.ToLower(sqlType) {
	case "int", "bigint", "smallint", "tinyint":
		return "int64"
	case "varchar", "text", "char":
		return "string"
	case "boolean":
		return "bool"
	default:
		return "interface{}" // 默认情况，使用空接口类型
	}
}

// GenerateModel 根据 TableInfo 生成 Go 代码文件
func GenerateModel(tableInfo TableInfo) {
	// Prepare file path
	filePath := filepath.Join("biz", "db", "model", strings.ToLower(tableInfo.TableName)+".go")
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	// Prepare template
	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Execute template
	err = tmpl.Execute(f, tableInfo)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	fmt.Printf("Model file generated: %s\n", filePath)
}

// SQLQueryInfo 包含 SQL 查询、API 路由信息和 SQL 解析后的 AST
type SQLQueryInfo struct {
	Name       string              // 方法名称，从 -- name: 后面的内容提取
	HTTPMethod string              // API HTTP 方法，从 --api.XXX: 后面的内容提取
	APIPath    string              // API 路由路径，从 --api.XXX: 后面的内容提取
	SQL        string              // SQL 查询语句，从 SQL 语句中提取
	ExecType   string              // 返回的数据量，例如 exec、many 等
	Table      string              // 操作的表名
	Params     []string            //需要的参数
	AST        sqlparser.Statement // SQL 解析后的抽象语法树
}

// ParseSQLQueryToStruct 解析 SQL 查询和 API 路由信息
func ParseSQLQueryToStruct(sql string) []SQLQueryInfo {
	// 正则表达式匹配 -- name: 和 --api.XXX:
	re := regexp.MustCompile(`--\s*name:\s*([^\n]+)\n--\s*api\.(?P<Method>[a-zA-Z]+):\s*([^\n]+)\n([^;]+);`)

	matches := re.FindAllStringSubmatch(sql, -1)
	var result []SQLQueryInfo

	for _, match := range matches {
		if len(match) < 5 {
			continue
		}

		fullName := strings.TrimSpace(match[1])
		method := strings.TrimSpace(match[2])
		apiPath := strings.TrimSpace(match[3])
		querySQL := strings.TrimSpace(match[4])
		execType := extractExecType(fullName)

		// 提取name中第一个:前的部分
		name := strings.Split(fullName, ":")[0]

		// 解析 SQL 查询语句为 AST
		stmt, err := sqlparser.Parse(querySQL)
		if err != nil {
			fmt.Printf("Failed to parse SQL: %v\n", err)
			continue
		}

		//// 提取表名
		//tableName := extractTableName(stmt)
		//todo:临时改名
		tableName := "user"

		info := SQLQueryInfo{
			Name:       name,
			HTTPMethod: method,
			APIPath:    apiPath,
			SQL:        querySQL,
			ExecType:   execType,
			Table:      tableName,
			AST:        stmt, // 存储解析后的语句，类型为 sqlparser.Statement
		}

		result = append(result, info)
	}

	return result
}

// extractExecType 从方法名称中提取执行类型（exec、many 等）
func extractExecType(name string) string {
	if strings.Contains(name, ":exec") {
		return "exec"
	} else if strings.Contains(name, ":many") {
		return "many"
	}
	// 默认为 exec 类型
	return "exec"
}

// extractTableName 从 AST 中提取表名
func extractTableName(stmt sqlparser.Statement) string {
	switch node := stmt.(type) {
	case *sqlparser.Insert:
		return node.Table.Name.String()
	case *sqlparser.Update:
		return node.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(*sqlparser.TableName).Name.String()
	case *sqlparser.Delete:
		return node.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(*sqlparser.TableName).Name.String()
	case *sqlparser.Select:
		return node.From[0].(*sqlparser.AliasedTableExpr).Expr.(*sqlparser.TableName).Name.String()
	default:
		return ""
	}
}

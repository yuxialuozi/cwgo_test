package anylize

import (
	"cwgo_test/codegen"
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func SQLExecutionShow(info []codegen.SQLQueryInfo, dsn string, name string) {
	// 查找与name匹配的SQL语句
	var query string
	found := false
	for _, q := range info {
		if q.Name == name {
			query = q.SQL
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Query named '%s' not found", name)
	}

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// 执行EXPLAIN语句
	explainQuery := fmt.Sprintf("EXPLAIN %s", query)
	rows, err := db.Query(explainQuery)
	if err != nil {
		log.Fatalf("Error executing EXPLAIN query: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error getting columns: %v", err)
	}

	// 创建一个slice of interface{}'s来保存每个列值
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 使用tabwriter来格式化输出
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	// 打印列名
	for _, col := range columns {
		fmt.Fprintf(w, "%s\t", col)
	}
	fmt.Fprintln(w)

	// 打印分隔行
	for range columns {
		fmt.Fprintf(w, "--------\t")
	}
	fmt.Fprintln(w)

	// 打印每一行数据
	for rows.Next() {
		// 将行数据保存到slice of interface{}'s
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// 打印行数据
		for _, col := range values {
			if col == nil {
				fmt.Fprintf(w, "NULL\t")
			} else {
				fmt.Fprintf(w, "%s\t", col)
			}
		}
		fmt.Fprintln(w)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error with rows: %v", err)
	}
}

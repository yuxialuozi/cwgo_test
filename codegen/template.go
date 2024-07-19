package codegen

const modelTemplate = `
package {{.PackageName}}

// {{.TableName}} represents the database table.
type {{.TableName}} struct {
{{- range .Columns}}
	{{.FieldName}} {{.GoType}} ` + "`json:\"{{.JsonTag}}\" query:\"{{.QueryTag}}\" body:\"{{.BodyTag}}\"`" + `
{{- end}}
}
`

const queryDbTemplate = `
package query

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
`

var templateContent = `
package query

import (
	"context"
	"database/sql"
	model "cwgo_test/biz/db/model"
)

{{range .Queries}}
const {{.Name}}Query = ` + "`{{.SQL}}`\n\n" + `

{{if eq .ExecType "exec"}}
// {{.Name}} executes the {{.Name}} query against the database and returns sql.Result.
func (q *Queries) {{.Name}}(ctx context.Context, {{.Table}} model.{{.Table}}) (sql.Result, error) {
	return q.db.ExecContext(ctx, {{.Name}}Query, {{.Table}})
}

{{else if eq .ExecType "execresult"}}
// {{.Name}} executes the {{.Name}} query against the database and returns an ExecResult.
func (q *Queries) {{.Name}}(ctx context.Context, {{.Table}} model.{{.Table}}) (ExecResult, error) {
	result, err := q.db.ExecContext(ctx, {{.Name}}Query, args...)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Result: result}, nil
}

{{else if eq .ExecType "one"}}
// {{.Name}} executes the {{.Name}} query against the database and returns a single row (*sql.Row).
func (q *Queries) {{.Name}}(ctx context.Context, {{.Table}} model.{{.Table}}) *sql.Row {
	return q.db.QueryRowContext(ctx, {{.Name}}Query, {{.Table}} )
}

{{else if eq .ExecType "many"}}
// {{.Name}} executes the {{.Name}} query against the database and returns multiple rows (*sql.Rows).
func (q *Queries) {{.Name}}(ctx context.Context, {{.Table}} model.{{.Table}}) (*sql.Rows, error) {
	return q.db.QueryContext(ctx, {{.Name}}Query, {{.Table}})
}

{{end}}
{{end}}
`

type query struct {
	Queries []SQLQueryInfo
	Table   TableInfo
}

const serverContent = `
package service

import (
    "context"
    "cwgo_test/biz/db/query"
    model "cwgo_test/biz/db/model"
    "log"
)

type {{.StructName}} struct {
    Query *query.Queries
}

{{range .Queries}}
func (s *{{$.StructName}}) {{.Name}}(ctx context.Context, {{.Table}} model.{{.Table}}) (interface{}, error) {
    result, err := s.Query.{{.Name}}(ctx, {{.Table}})
    if err != nil {
        log.Printf("Error executing {{.Name}}: %v", err)
        return nil, err
    }
    return result, nil
}
{{end}}
`
const routerContent = `
package router

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "cwgo_test/biz/service"
    model "cwgo_test/biz/db/model"
)

func RegisterRoutes(h *server.Hertz) {
    {{range .Queries}}
    h.{{.HTTPMethod}}("{{.APIPath}}", func(c context.Context, ctx *app.RequestContext) {
        var user model.User
        if err := ctx.Bind(&user); err != nil {
            ctx.JSON(400, map[string]string{"error": "bad request"})
            return
        }

        result, err := service.{{.Name}}(c, &user)
        if err != nil {
            ctx.JSON(500, map[string]string{"error": "internal server error"})
            return
        }

        ctx.JSON(200, result)
    })
    {{end}}
}
`
const mainContent = `
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
    "github.com/cloudwego/hertz/pkg/app/server"
    "cwgo_test/biz/router" // Assume this is the generated router package
    "cwgo_test/biz/service" // Assuming this is where your services are defined
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

`

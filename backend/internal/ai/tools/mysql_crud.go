package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlCrudInput struct {
	DSN         string `json:"dsn" jsonschema:"description=The Data Source Name for connecting to the MySQL database, including username, password, host, port, and database name"`
	SQL         string `json:"sql" jsonschema:"description=The SQL query to execute against the MySQL database"`
	OperateType string `json:"operate_type" jsonschema:"description=The type of SQL operation to perform: query, insert, update, or delete"`
}

func NewMysqlCrudTool() (tool.InvokableTool, error) {
	return utils.InferOptionableTool(
		"mysql_crud",
		"Execute SQL queries against the MySQL database and return results in JSON format. Use this tool when you need to query, insert, update or delete data from the database. The results will be formatted as JSON for easy parsing.",
		func(ctx context.Context, input *MysqlCrudInput, opts ...tool.Option) (output string, err error) {
			db, err := gorm.Open(mysql.Open(input.DSN), &gorm.Config{})
			if err != nil {
				log.Printf("failed to connect to MySQL: %v", err)
				return "", err
			}

			err = db.Exec(input.SQL).Error
			if err != nil {
				log.Printf("failed to execute SQL: %v", err)
				return "", err
			}
			if input.OperateType == "query" {
				var results []interface{}
				err = db.Raw(input.SQL).Scan(&results).Error
				if err != nil {
					log.Printf("failed to scan results: %v", err)
					return "", err
				}
				resBytes, err := json.Marshal(results)
				return string(resBytes), err
			}
			return "", nil
		})
}

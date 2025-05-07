package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/robertgouveia/do-my-job/storage"
)

func Connect(serverName string) (*sql.DB, error) {
	s, err := storage.LoadServerConfig(serverName)
	if err != nil {
		return nil, fmt.Errorf("server config error: %v", err)
	}

	connStr := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;database=%s;",
		s.Host, s.Username, s.Password, s.Database,
	)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v -- conn: %s", err, connStr)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open database: %v -- conn: %s", err, connStr)
	}

	return db, nil
}

func Execute(db *sql.DB, stmt string, params ...interface{}) (sql.Result, string, error) {
	paramDebug := ""
	for i, param := range params {
		paramDebug += fmt.Sprintf("[%d : %v] ", i, param)
	}

	preparedStmt, err := db.Prepare(stmt)
	if err != nil {
		return nil, paramDebug, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer preparedStmt.Close()

	result, err := preparedStmt.Exec(params...)
	if err != nil {
		return nil, paramDebug, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, paramDebug, nil
}

func CreateNamedParameters(values ...interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	for i, value := range values {
		params[fmt.Sprintf("p%d", i+1)] = value
	}
	return params
}

func ExecuteWithNamedParams(db *sql.DB, stmt string, params map[string]interface{}) (sql.Result, string, error) {
	paramDebug := ""
	for name, value := range params {
		paramDebug += fmt.Sprintf("[%s : %v] ", name, value)
	}

	preparedStmt, err := db.Prepare(stmt)
	if err != nil {
		return nil, paramDebug, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer preparedStmt.Close()

	var orderedParams []interface{}
	paramNames := extractParamNames(stmt)
	for _, name := range paramNames {
		if value, exists := params[name]; exists {
			orderedParams = append(orderedParams, sql.Named(name, value))
		} else {
			return nil, paramDebug, fmt.Errorf("missing parameter: %s", name)
		}
	}

	result, err := preparedStmt.Exec(orderedParams...)
	if err != nil {
		return nil, paramDebug, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, paramDebug, nil
}

func extractParamNames(stmt string) []string {
	var params []string
	parts := strings.Split(stmt, "@")
	for i := 1; i < len(parts); i++ {
		var name string
		for j := 0; j < len(parts[i]); j++ {
			if (parts[i][j] >= 'a' && parts[i][j] <= 'z') ||
				(parts[i][j] >= 'A' && parts[i][j] <= 'Z') ||
				(parts[i][j] >= '0' && parts[i][j] <= '9') ||
				parts[i][j] == '_' {
				name += string(parts[i][j])
			} else {
				break
			}
		}
		params = append(params, name)
	}
	return params
}

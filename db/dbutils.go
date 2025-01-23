package dbutils

import (
	"database/sql"
	"fmt"
	"strings"
	"testcaseautomation/constants"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectToSQLServer() (*sql.DB, error) {
	connString := "server=wva60bddpadb11v.corp.bankofamerica.com;port=15001;trusted_connection=yes;user id=CORP\\zsppuwid;password=c5prrP4s;database=idwdb"
	conn, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func UpdateTable(db *sql.DB, tableName string, whereCondition string, updateValues map[string]interface{}) error {
	query := "UPDATE " + tableName + " SET "
	values := make([]interface{}, 0)
	for key, value := range updateValues {
		query += key + " = ?, "
		values = append(values, value)
	}
	query = query[:len(query)-2] // to remove the last comma and space
	query += " WHERE " + whereCondition

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFromTable(db *sql.DB, tableName string, whereCondition string, args ...interface{}) error {
	query := "DELETE FROM " + tableName + " WHERE " + whereCondition

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func InsertScanWindowRecords(db *sql.DB, aitNo, dbType, startTime, endTime string) error {
	// Delete existing records for the specified AIT_NO
	_, err := db.Exec("DELETE FROM AIT_SCAN_WINDOW WHERE AIT_NO = ?", aitNo)
	if err != nil {
		return err
	}

	windowStartTime := "00:00:00"
	if startTime != "" {
		windowStartTime = startTime
	}

	windowEndTime := "23:59:00"
	if endTime != "" {
		windowEndTime = endTime
	}

	lastUpdatedUser := "testframework"
	profile := "public"

	// Insert new records
	for _, day := range constants.Days {
		currentTime := time.Now().Format("2006-01-02 15:04:05") // Get current timestamp
		_, err := db.Exec("INSERT INTO AIT_SCAN_WINDOW (AIT_NO, DB_TYPE, SCAN_DAY, START_TIME_EST, END_TIME_EST, LAST_UPDATED, LAST_UPDATED_USER, PROFILE) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			aitNo, dbType, day, windowStartTime, windowEndTime, currentTime, lastUpdatedUser, profile)
		if err != nil {
			return err
		}
	}
	return nil
}

func SelectAndInsertRowsToTable(db *sql.DB, selectQuery string, newValues map[string]interface{}) error {
	// Execute the SELECT query
	rows, err := db.Query(selectQuery)
	if err != nil {
		return fmt.Errorf("error executing select query: %v", err)
	}
	defer rows.Close()

	// Get column names from the SELECT query
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}

	// Prepare a slice to hold the column values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Iterate over the rows and convert to INSERT statements
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Update values with new values
		for col, newVal := range newValues {
			for i, column := range columns {
				if column == col {
					values[i] = newVal
				}
			}
		}

		// Build the INSERT statement
		insertQuery := buildInsertQuery("KAFKA_STAT", columns, values)
		_, err = db.Exec(insertQuery)
		if err != nil {
			return fmt.Errorf("error executing insert query: %v", err)
		}
	}

	return nil
}

func buildInsertQuery(table string, columns []string, values []interface{}) string {
	columnsStr := strings.Join(columns, ", ")
	valuesStr := make([]string, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case string:
			valuesStr[i] = fmt.Sprintf("'%s'", v)
		case []byte:
			valuesStr[i] = fmt.Sprintf("'%s'", string(v))
		default:
			valuesStr[i] = fmt.Sprintf("%v", v)
		}
	}
	valuesJoined := strings.Join(valuesStr, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columnsStr, valuesJoined)
}
